package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	COMPLETIONURL = "https://api.openai.com/v1/chat/completions"
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

type openaiUsage struct {
	PromptTokens            int            `json:"prompt_tokens"`
	CompletionTokens        int            `json:"completion_tokens"`
	TotalTokens             int            `json:"total_tokens"`
	CompletionTokensDetails map[string]any `json:"completion_tokens_details"`
}

type openaiChoice struct {
	Index   int     `json:"index"`
	Message message `json:"message"`
}

type openaiError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type openaiClient struct {
	apiKey string
	// add streaming channel?
}

func NewOpenaiClient(apiKey string) (openaiClient, error) {
	if apiKey == "" {
		return openaiClient{}, errors.New("Missing OpenAI API key.")
	}
	return openaiClient{apiKey: apiKey}, nil
}

func (oc openaiClient) Complete(options ...completionOption) (openaiRequest, openaiResponse, error) {

	request, err := NewOpenaiRequest(options...)
	if err != nil {
		return openaiRequest{}, openaiResponse{}, err
	}

	res, err := makeHTTPRequest(request, oc)
	if err != nil {
		return openaiRequest{}, openaiResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return openaiRequest{}, openaiResponse{}, err
	}

	openaiRes := new(openaiResponse)
	json.Unmarshal(body, openaiRes)

	// attach status code to response object
	openaiRes.StatusCode = res.StatusCode

	return *request, *openaiRes, nil
}

func (oc openaiClient) StreamComplete(options ...completionOption) (openaiRequest, openaiResponse, error) {

	request, err := NewOpenaiRequest(options...)

	if err != nil {
		return openaiRequest{}, openaiResponse{}, err
	}
	request.Stream = true

	res, err := makeHTTPRequest(request, oc)
	if err != nil {
		return openaiRequest{}, openaiResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return openaiRequest{}, openaiResponse{}, err
	}

	openaiRes := new(openaiResponse)
	json.Unmarshal(body, openaiRes)

	// attach status code to response object
	openaiRes.StatusCode = res.StatusCode

	return *request, *openaiRes, nil
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

	req, err := http.NewRequest(http.MethodPost, COMPLETIONURL, bytes.NewReader(jsonRequest))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", oc.apiKey))

	client := http.Client{Timeout: time.Duration(30 * time.Second)}
	res, err := client.Do(req)

	return res, err
}
