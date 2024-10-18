package client

import (
	"github.com/azr4e1/gocast/client/openai"
)

type llmProvider int
type clientOption func(*client) error

type LLMClient interface {
	Complete([]Message) Response
}

const (
	OPENAI llmProvider = iota
	OLLAMA
)

type client struct {
	provider      llmProvider
	model         string
	api_key       string
	api_base      string
	backendClient LLMClient
}

func NewClient(...clientOption) (client, error)

func (c client) Complete()
