package gollum

import (
	ll "github.com/azr4e1/gollum/ollama"
	oai "github.com/azr4e1/gollum/openai"
)

func (c llmClient) ToOpenAI() (oai.OpenaiClient, error) {
	client, err := oai.NewClient(c.apiKey)
	if err != nil {
		return oai.OpenaiClient{}, err
	}
	client.Timeout = c.Timeout

	return client, nil
}

func (c llmClient) ToOllama() (ll.OllamaClient, error) {
	client, err := ll.NewClient(c.apiBase)
	if err != nil {
		return ll.OllamaClient{}, err
	}
	client.Timeout = c.Timeout

	return client, nil
}
