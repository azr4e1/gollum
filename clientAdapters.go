package gollum

import (
	oai "github.com/azr4e1/gollum/openai"
)

func (c llmClient) ToOpenAI() (oai.OpenaiClient, error) {
	client, err := oai.NewClient(c.apiKey)
	if err != nil {
		return oai.OpenaiClient{}, err
	}

	return client, nil
}
