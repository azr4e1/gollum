package gollum

import (
	"context"
	"errors"
	"fmt"

	m "github.com/azr4e1/gollum/message"
)

type CompletionType string

const (
	Text     CompletionType = "text"
	ToolCall CompletionType = "tool"
)

type CompletionRequest struct {
	Model               string          `json:"model"`
	System              m.Message       `json:"system_message"`
	Messages            []m.Message     `json:"messages"`
	Tools               []Tool          `json:"tools,omitempty"`
	Stream              bool            `json:"stream"`
	FreqPenalty         *float64        `json:"frequency_penalty,omitempty"`
	LogitBias           map[int]int     `json:"logit_bias,omitempty"`
	LogProbs            *bool           `json:"logprobs,omitempty"`
	TopLogProbs         *int            `json:"top_logprobs,omitempty"`
	MaxCompletionTokens *int            `json:"max_tokens,omitempty"`
	PresencePenalty     *float64        `json:"presence_penalty,omitempty"`
	Seed                *int            `json:"seed,omitempty"`
	Stop                []string        `json:"stop,omitempty"`
	Temperature         *float64        `json:"temperature,omitempty"`
	TopP                *float64        `json:"top_p,omitempty"`
	TopK                *int            `json:"top_k,omitempty"`
	User                string          `json:"user,omitempty"`
	Ctx                 context.Context `json:"-"`
}

type CompletionUsage struct {
	PromptTokens            int            `json:"prompt_tokens"`
	CompletionTokens        int            `json:"completion_tokens"`
	TotalTokens             int            `json:"total_tokens"`
	CompletionTokensDetails map[string]any `json:"completion_tokens_details"`
}

type CompletionError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type CompletionResponse struct {
	Id         string          `json:"id"`
	Object     string          `json:"object"`
	Created    int             `json:"created"`
	Model      string          `json:"model"`
	Type       CompletionType  `json:"type"`
	Message    m.Message       `json:"message"`
	Done       bool            `json:"done"`
	Usage      CompletionUsage `json:"usage"`
	Error      CompletionError `json:"error,omitempty"`
	StatusCode int             `json:"status_code"`
}

func (or CompletionResponse) Content() string {
	return or.Message.Content
}

func (or CompletionResponse) Tools() []m.ToolCall {
	return or.Message.ToolCalls
}

func (or CompletionResponse) Tool() (m.ToolCall, error) {
	if tools := or.Tools(); tools != nil && len(tools) > 0 {
		return tools[0], nil
	}
	return m.ToolCall{}, errors.New("No tools available.")
}

func (or CompletionResponse) Err() error {
	if or.Error.Type == "" && or.Error.Message == "" {
		return nil
	}
	return errors.New(fmt.Sprintf("%s: %s", or.Error.Type, or.Error.Message))
}

func NewCompletionRequest(options ...completionOption) (*CompletionRequest, error) {
	request := new(CompletionRequest)

	for _, o := range options {
		err := o(request)
		if err != nil {
			return &CompletionRequest{}, err
		}
	}

	if request.Model == "" {
		return &CompletionRequest{}, errors.New("Missing model name.")
	}
	if m := request.Messages; m == nil || len(m) == 0 {
		return &CompletionRequest{}, errors.New("Missing messages to send.")
	}

	return request, nil
}
