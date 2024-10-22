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

type openaiClient struct {
	apiKey        string
	streamChannel chan completionResponse
}

func NewClient(apiKey string) (openaiClient, error) {
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

func (oc *openaiClient) EnableStream() <-chan completionResponse {
	c := make(chan completionResponse)
	oc.streamChannel = c

	return c
}

func (oc openaiClient) IsStreaming() bool {
	return oc.streamChannel == nil
}

func (oc openaiClient) Complete(options ...completionOption) (completionRequest, completionResponse, error) {
	request, err := NewCompletionRequest(options...)
	if err != nil {
		return *request, completionResponse{}, err
	}
	if oc.streamChannel != nil {
		request.Stream = true
	}

	res, err := makeHTTPCompletionRequest(request, oc)
	if err != nil {
		return *request, completionResponse{}, err
	}
	defer res.Body.Close()

	if oc.streamChannel != nil {
		err = oc.readCompletionStreamResponse(res)
		return *request, completionResponse{}, err
	}

	openaiRes, err := oc.readCompletionResponse(res)
	return *request, openaiRes, err
}

func (oc openaiClient) readCompletionResponse(res *http.Response) (completionResponse, error) {

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return completionResponse{}, err
	}

	openaiRes := new(completionResponse)
	err = json.Unmarshal(body, openaiRes)
	if err != nil {
		return completionResponse{}, err
	}

	// attach status code to response object
	openaiRes.StatusCode = res.StatusCode

	return *openaiRes, nil
}

func (oc openaiClient) readCompletionStreamResponse(res *http.Response) error {

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

		chunk := new(completionResponse)
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

	openaiRes := new(completionResponse)
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

func NewCompletionRequest(options ...completionOption) (*completionRequest, error) {
	request := new(completionRequest)

	for _, o := range options {
		err := o(request)
		if err != nil {
			return &completionRequest{}, err
		}
	}

	if request.Model == "" {
		return &completionRequest{}, errors.New("Missing model name.")
	}
	if m := request.Messages; m == nil || len(m) == 0 {
		return &completionRequest{}, errors.New("Missing messages to send.")
	}

	return request, nil
}

func makeHTTPCompletionRequest(request *completionRequest, oc openaiClient) (*http.Response, error) {
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
