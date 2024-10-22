package gollum

import (
	"errors"
)

type llmProvider int

const (
	OPENAI llmProvider = iota + 1
	OLLAMA
)

type llmClient struct {
	provider llmProvider
	apiKey   string
	apiBase  string
}

func NewClient(options ...clientOption) (llmClient, error) {
	client := new(llmClient)
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

func (c llmClient) Complete(options ...completionOption) (CompletionRequest, CompletionResponse, error) {
	request, err := NewCompletionRequest(options...)
	if err != nil {
		return *request, CompletionResponse{}, err
	}

	switch c.provider {
	case OPENAI:
		openaiReq := request.ToOpenAI()
		openaiClient, err := c.ToOpenAI()
		if err != nil {
			return *request, CompletionResponse{}, err
		}
		_, result, err := openaiClient.CompleteWithCustomRequest(&openaiReq)
		if err != nil {
			return *request, CompletionResponse{}, err
		}

		return *request, ResponseFromOpenAI(result), nil
	}

	return *request, CompletionResponse{}, nil
}

// func (c llmClient) Speech()
