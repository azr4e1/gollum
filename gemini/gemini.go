package gemini

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
	completionURL = "https://generativelanguage.googleapis.com/v1beta/models/%s:%s?alt=sse&key=%s"
)

const (
	streamEnd  = "data: [DONE]"
	dataPrefix = "data: "
)

type StreamingFunction func(CompletionResponse) error

type GeminiClient struct {
	apiKey         string
	stream         bool
	streamFunction StreamingFunction
	Timeout        time.Duration
}

func NewClient(apiKey string) (GeminiClient, error) {
	if apiKey == "" {
		return GeminiClient{}, errors.New("Missing Gemini API key.")
	}
	return GeminiClient{apiKey: apiKey, Timeout: 30 * time.Second}, nil
}

func (oc *GeminiClient) EnableStream(function StreamingFunction) {
	oc.stream = true
	oc.streamFunction = function
}

func (oc GeminiClient) Complete(request *CompletionRequest) (CompletionRequest, CompletionResponse, error) {
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

	geminiRes, err := oc.readCompletionResponse(res)
	return *request, geminiRes, err
}

func (oc GeminiClient) readCompletionResponse(res *http.Response) (CompletionResponse, error) {

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return CompletionResponse{}, err
	}

	// remove data prefix from response
	if string(body)[:len(dataPrefix)] == dataPrefix {
		body = body[len([]byte(dataPrefix)):]
	}

	geminiRes := new(CompletionResponse)
	err = json.Unmarshal(body, geminiRes)
	if err != nil {
		return CompletionResponse{}, err
	}

	// attach status code to response object
	geminiRes.StatusCode = res.StatusCode

	return *geminiRes, geminiRes.err()
}

func (oc GeminiClient) readCompletionStreamResponse(res *http.Response) error {

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

	geminiRes := new(CompletionResponse)
	err = json.Unmarshal(body, geminiRes)
	if err != nil {
		return err
	}

	// attach status code to response object
	geminiRes.StatusCode = res.StatusCode

	return geminiRes.err()
}
