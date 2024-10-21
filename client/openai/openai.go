package openai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	completionURL = "https://api.openai.com/v1/chat/completions"
)

const (
	streamEnd  = "data: [DONE]"
	dataPrefix = "data: "
)

type openaiRequest struct {
	Model               string      `json:"model"`
	Messages            []message   `json:"messages"`
	Stream              bool        `json:"stream"`
	FreqPenalty         *float64    `json:"frequency_penalty,omitempty"`
	LogitBias           map[int]int `json:"logit_bias,omitempty"`
	LogProbs            *bool       `json:"logprobs,omitempty"`
	TopLogProbs         *int        `json:"top_logprobs,omitempty"`
	MaxCompletionTokens *int        `json:"max_completion_tokens,omitempty"`
	CompletionChoices   *int        `json:"n,omitempty"`
	PresencePenalty     *float64    `json:"presence_penalty,omitempty"`
	Seed                *int        `json:"seed,omitempty"`
	Stop                []string    `json:"stop,omitempty"`
	Temperature         *float64    `json:"temperature,omitempty"`
	TopP                *float64    `json:"top_p,omitempty"`
	User                string      `json:"user,omitempty"`
}

type openaiUsage struct {
	PromptTokens            int            `json:"prompt_tokens"`
	CompletionTokens        int            `json:"completion_tokens"`
	TotalTokens             int            `json:"total_tokens"`
	CompletionTokensDetails map[string]any `json:"completion_tokens_details"`
}

type openaiChoice struct {
	Index   int     `json:"index"`
	Message message `json:"message,omitempty"`
	Delta   message `json:"delta,omitempty"`
}

type openaiError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type openaiResponse struct {
	Id         string         `json:"id"`
	Object     string         `json:"object"`
	Created    int            `json:"created"`
	Model      string         `json:"model"`
	Choices    []openaiChoice `json:"choices"`
	Usage      openaiUsage    `json:"usage"`
	Error      openaiError    `json:"error"`
	StatusCode int            `json:"status_code"`
}

func (or openaiResponse) IsEOS() bool {
	return or.Error.Type == streamEnd
}

func newEOS(message string) openaiResponse {
	return openaiResponse{Error: openaiError{Type: streamEnd, Message: message}}
}

type openaiClient struct {
	apiKey        string
	streamChannel chan openaiResponse
}

func NewOpenaiClient(apiKey string) (openaiClient, error) {
	if apiKey == "" {
		return openaiClient{}, errors.New("Missing OpenAI API key.")
	}
	return openaiClient{apiKey: apiKey}, nil
}

func (oc *openaiClient) DisableStream() {
	if c := oc.streamChannel; c != nil {
		close(c)
	}
	oc.streamChannel = nil
}

func (oc *openaiClient) EnableStream() <-chan openaiResponse {
	c := make(chan openaiResponse)
	oc.streamChannel = c

	return c
}

func (oc openaiClient) IsStreaming() bool {
	return oc.streamChannel == nil
}

func (oc openaiClient) Complete(options ...completionOption) (openaiRequest, openaiResponse, error) {
	request, err := NewOpenaiRequest(options...)
	if err != nil {
		return *request, openaiResponse{}, err
	}
	if oc.streamChannel != nil {
		request.Stream = true
	}

	res, err := makeHTTPRequest(request, oc)
	if err != nil {
		return *request, openaiResponse{}, err
	}
	defer res.Body.Close()

	if oc.streamChannel != nil {
		err = oc.readStreamResponse(res)
		return *request, openaiResponse{}, err
	}

	openaiRes, err := oc.readResponse(res)
	return *request, openaiRes, err
}

func (oc openaiClient) readResponse(res *http.Response) (openaiResponse, error) {

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return openaiResponse{}, err
	}

	openaiRes := new(openaiResponse)
	err = json.Unmarshal(body, openaiRes)
	if err != nil {
		return openaiResponse{}, err
	}

	// attach status code to response object
	openaiRes.StatusCode = res.StatusCode

	return *openaiRes, nil
}

func (oc openaiClient) readStreamResponse(res *http.Response) error {

	reader := bufio.NewReader(res.Body)

	// read response body until end of stream
	for res.StatusCode == http.StatusOK {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				oc.streamChannel <- newEOS("End of stream.")
				return nil
			}
			oc.streamChannel <- newEOS("Byte read error.")
			return err
		}

		line = bytes.TrimSpace(line)
		// skip blank lines
		if len(line) == 0 {
			continue
		}

		if string(line) == streamEnd {
			oc.streamChannel <- newEOS("End of stream.")
			return nil
		}

		// remove data prefix from response
		if string(line)[:len(dataPrefix)] == dataPrefix {
			line = line[len([]byte(dataPrefix)):]
		}

		chunk := new(openaiResponse)
		err = json.Unmarshal(line, chunk)
		if err != nil {
			oc.streamChannel <- newEOS("JSON unmarshalling error.")
			return err
		}
		// attach status code to response object
		chunk.StatusCode = res.StatusCode

		oc.streamChannel <- *chunk
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		oc.streamChannel <- newEOS("Byte read error.")
		return err
	}

	openaiRes := new(openaiResponse)
	err = json.Unmarshal(body, openaiRes)
	if err != nil {
		oc.streamChannel <- newEOS("JSON unmarshalling error.")
		return err
	}

	// attach status code to response object
	openaiRes.StatusCode = res.StatusCode
	oc.streamChannel <- *openaiRes

	return nil
}

func NewOpenaiRequest(options ...completionOption) (*openaiRequest, error) {
	request := new(openaiRequest)

	for _, o := range options {
		err := o(request)
		if err != nil {
			return &openaiRequest{}, err
		}
	}

	if request.Model == "" {
		return &openaiRequest{}, errors.New("Missing model name.")
	}
	if m := request.Messages; m == nil || len(m) == 0 {
		return &openaiRequest{}, errors.New("Missing messages to send.")
	}

	return request, nil
}

func makeHTTPRequest(request *openaiRequest, oc openaiClient) (*http.Response, error) {
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, completionURL, bytes.NewReader(jsonRequest))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", oc.apiKey))

	client := http.Client{Timeout: time.Duration(30 * time.Second)}
	res, err := client.Do(req)

	return res, err
}
