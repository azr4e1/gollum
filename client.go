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

type LLMClient struct {
	provider       llmProvider
	apiKey         string
	apiBase        string
	stream         bool
	streamFunction StreamingFunction
	Timeout        time.Duration
}

func NewClient(options ...clientOption) (LLMClient, error) {
	client := new(LLMClient)
	client.Timeout = 30 * time.Second
	for _, o := range options {
		err := o(client)
		if err != nil {
			return LLMClient{}, err
		}
	}

	if client.provider == 0 {
		return LLMClient{}, errors.New("provider is empty.")
	}
	if client.apiKey == "" && client.apiBase == "" {
		return LLMClient{}, errors.New("must provide at least one of apiKey or apiBase")
	}

	return *client, nil
}

func (oc *LLMClient) DisableStream() {
	oc.stream = false
	oc.streamFunction = nil
}

func (oc *LLMClient) EnableStream(function StreamingFunction) {
	oc.stream = true
	oc.streamFunction = function
}

func (oc LLMClient) IsStreaming() bool {
	return oc.stream
}

func (c LLMClient) Complete(options ...completionOption) (CompletionRequest, CompletionResponse, error) {
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

func (c LLMClient) TextToSpeech(options ...speechOption) (TTSRequest, TTSResponse, error) {
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
