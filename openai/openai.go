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

func (oc openaiClient) Speech(opts ...audioOption) (audioRequest, audioResponse, error) {
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
