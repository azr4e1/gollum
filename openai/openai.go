package openai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

const (
	completionURL = "https://api.openai.com/v1/chat/completions"
	speechURL     = "https://api.openai.com/v1/audio/speech"
)

const (
	streamEnd  = "data: [DONE]"
	dataPrefix = "data: "
)

type StreamingFunction func(CompletionResponse) error

type OpenaiClient struct {
	apiKey         string
	stream         bool
	streamFunction StreamingFunction
	Timeout        time.Duration
}

func NewClient(apiKey string) (OpenaiClient, error) {
	if apiKey == "" {
		return OpenaiClient{}, errors.New("Missing OpenAI API key.")
	}
	return OpenaiClient{apiKey: apiKey, Timeout: 30 * time.Second}, nil
}

func (oc *OpenaiClient) EnableStream(function StreamingFunction) {
	oc.stream = true
	oc.streamFunction = function
}

func (oc OpenaiClient) Complete(request *CompletionRequest) (CompletionRequest, CompletionResponse, error) {
	request.Stream = oc.stream

	res, err := makeHTTPCompletionRequest(request, oc)
	if err != nil {
		return *request, CompletionResponse{}, err
	}
	defer res.Body.Close()

	if oc.stream {
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

	return *openaiRes, openaiRes.err()
}

func (oc OpenaiClient) readCompletionStreamResponse(res *http.Response) error {

	reader := bufio.NewReader(res.Body)

	// read response body until end of stream
	for res.StatusCode == http.StatusOK {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		line = bytes.TrimSpace(line)
		// skip blank lines
		if len(line) == 0 {
			continue
		}

		if string(line) == streamEnd {
			return nil
		}

		// remove data prefix from response
		if string(line)[:len(dataPrefix)] == dataPrefix {
			line = line[len([]byte(dataPrefix)):]
		}

		chunk := new(CompletionResponse)
		err = json.Unmarshal(line, chunk)
		if err != nil {
			return err
		}
		// attach status code to response object
		chunk.StatusCode = res.StatusCode

		err = oc.streamFunction(*chunk)
		if err != nil {
			return err
		}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	openaiRes := new(CompletionResponse)
	err = json.Unmarshal(body, openaiRes)
	if err != nil {
		return err
	}

	// attach status code to response object
	openaiRes.StatusCode = res.StatusCode

	return openaiRes.err()
}

func (oc OpenaiClient) TextToSpeech(request *TTSRequest) (TTSRequest, TTSResponse, error) {
	response := new(TTSResponse)

	res, err := makeHTTPTTSRequest(request, oc)
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
		return *request, *response, response.Err()
	}

	response.StatusCode = res.StatusCode
	response.Audio = body
	return *request, *response, response.Err()
}
