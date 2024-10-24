package gollum

import (
	"errors"
	"time"
)

type llmProvider int

const (
	OPENAI llmProvider = iota + 1
	OLLAMA
)

type StreamingFunction func(CompletionResponse) error

type llmClient struct {
	provider       llmProvider
	apiKey         string
	apiBase        string
	stream         bool
	streamFunction StreamingFunction
	Timeout        time.Duration
}

func NewClient(options ...clientOption) (llmClient, error) {
	client := new(llmClient)
	client.Timeout = 30 * time.Second
	for _, o := range options {
		err := o(client)
		if err != nil {
			return llmClient{}, nil
		}
	}

	if client.provider == 0 {
		return llmClient{}, errors.New("provider is empty.")
	}
	if client.apiKey == "" && client.apiBase == "" {
		return llmClient{}, errors.New("must provide at least one of apiKey or apiBase")
	}

	return *client, nil
}

func (oc *llmClient) DisableStream() {
	oc.stream = false
	oc.streamFunction = nil
}

func (oc *llmClient) EnableStream(function StreamingFunction) {
	oc.stream = true
	oc.streamFunction = function
}

func (oc llmClient) IsStreaming() bool {
	return oc.stream
}

func (c llmClient) Complete(options ...completionOption) (CompletionRequest, CompletionResponse, error) {
	request, err := NewCompletionRequest(options...)
	if err != nil {
		return *request, CompletionResponse{}, err
	}

	switch c.provider {
	case OPENAI:
		return openaiComplete(request, c)
	case OLLAMA:
		return ollamaComplete(request, c)
	}

	return *request, CompletionResponse{}, errors.New("completion not implemented for this provider.")
}

func (c llmClient) TextToSpeech(options ...speechOption) (TTSRequest, TTSResponse, error) {
	request, err := NewTTSRequest(options...)
	if err != nil {
		return *request, TTSResponse{}, err
	}

	switch c.provider {
	case OPENAI:
		return openaiTTS(request, c)
	}

	return *request, TTSResponse{}, errors.New("text to speech not implemented for this provider.")
}
