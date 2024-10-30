package ollama

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
	completionURL = "api/chat"
)

const (
	streamEnd  = "data: [DONE]"
	dataPrefix = "data: "
)

type StreamingFunction func(CompletionResponse) error

type OllamaClient struct {
	baseURL        string
	stream         bool
	streamFunction StreamingFunction
	Timeout        time.Duration
}

func NewClient(baseURL string) (OllamaClient, error) {
	if baseURL == "" {
		return OllamaClient{}, errors.New("Missing base URL.")
	}
	return OllamaClient{baseURL: baseURL, stream: false, Timeout: 30 * time.Second}, nil
}

func (oc *OllamaClient) EnableStream(function StreamingFunction) {
	oc.stream = true
	oc.streamFunction = function
}

func (oc OllamaClient) Complete(request *CompletionRequest) (CompletionRequest, CompletionResponse, error) {
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

	ollamaRes, err := oc.readCompletionResponse(res)
	return *request, ollamaRes, err
}

func (oc OllamaClient) readCompletionResponse(res *http.Response) (CompletionResponse, error) {

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return CompletionResponse{}, err
	}

	ollamaRes := new(CompletionResponse)
	err = json.Unmarshal(body, ollamaRes)
	if err != nil {
		return CompletionResponse{}, err
	}

	// attach status code to response object
	ollamaRes.StatusCode = res.StatusCode

	return *ollamaRes, ollamaRes.err()
}

func (oc OllamaClient) readCompletionStreamResponse(res *http.Response) error {

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

	ollamaRes := new(CompletionResponse)
	err = json.Unmarshal(body, ollamaRes)
	if err != nil {
		return err
	}

	// attach status code to response object
	ollamaRes.StatusCode = res.StatusCode

	return ollamaRes.err()
}
