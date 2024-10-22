package openai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const (
	completionURL = "https://api.openai.com/v1/chat/completions"
	audioURL      = "https://api.openai.com/v1/audio/speech"
)

const (
	streamEnd  = "data: [DONE]"
	dataPrefix = "data: "
)

type OpenaiClient struct {
	apiKey        string
	streamChannel chan CompletionResponse
}

func NewClient(apiKey string) (OpenaiClient, error) {
	if apiKey == "" {
		return OpenaiClient{}, errors.New("Missing OpenAI API key.")
	}
	return OpenaiClient{apiKey: apiKey}, nil
}

func (oc *OpenaiClient) DisableStream() {
	if c := oc.streamChannel; c != nil {
		close(c)
	}
	oc.streamChannel = nil
}

func (oc *OpenaiClient) EnableStream() <-chan CompletionResponse {
	c := make(chan CompletionResponse)
	oc.streamChannel = c

	return c
}

func (oc OpenaiClient) IsStreaming() bool {
	return oc.streamChannel == nil
}

func (oc OpenaiClient) Complete(options ...completionOption) (CompletionRequest, CompletionResponse, error) {
	request, err := NewCompletionRequest(options...)
	if err != nil {
		return *request, CompletionResponse{}, err
	}

	return oc.CompleteWithCustomRequest(request)
}

func (oc OpenaiClient) CompleteWithCustomRequest(request *CompletionRequest) (CompletionRequest, CompletionResponse, error) {
	if oc.streamChannel != nil {
		request.Stream = true
	}

	res, err := makeHTTPCompletionRequest(request, oc)
	if err != nil {
		return *request, CompletionResponse{}, err
	}
	defer res.Body.Close()

	if oc.streamChannel != nil {
		err = oc.readCompletionStreamResponse(res)
		return *request, CompletionResponse{}, err
	}

	openaiRes, err := oc.readCompletionResponse(res)
	return *request, openaiRes, err
}

func (oc OpenaiClient) readCompletionResponse(res *http.Response) (CompletionResponse, error) {

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return CompletionResponse{}, err
	}

	openaiRes := new(CompletionResponse)
	err = json.Unmarshal(body, openaiRes)
	if err != nil {
		return CompletionResponse{}, err
	}

	// attach status code to response object
	openaiRes.StatusCode = res.StatusCode

	return *openaiRes, nil
}

func (oc OpenaiClient) readCompletionStreamResponse(res *http.Response) error {

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

		chunk := new(CompletionResponse)
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

	openaiRes := new(CompletionResponse)
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

func (oc OpenaiClient) Speech(opts ...audioOption) (audioRequest, audioResponse, error) {
	request, err := NewAudioRequest(opts...)
	response := new(audioResponse)
	if err != nil {
		return *request, *response, err
	}

	res, err := makeHTTPAudioRequest(request, oc)
	if err != nil {
		return *request, *response, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return *request, *response, err
	}

	// check the return status
	if res.StatusCode != http.StatusOK {
		err := json.Unmarshal(body, response)
		if err != nil {
			return *request, *response, err
		}

		response.StatusCode = res.StatusCode
		return *request, *response, err
	}

	response.StatusCode = res.StatusCode
	response.Audio = body
	return *request, *response, nil
}
