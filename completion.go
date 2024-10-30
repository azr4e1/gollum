package gollum

import (
	"context"
	"errors"
	"fmt"
	m "github.com/azr4e1/gollum/message"
)

const streamEnd = "gollum_end_of_stream"

type CompletionRequest struct {
	Model               string
	Messages            []m.Message
	Stream              bool
	FreqPenalty         *float64
	LogitBias           map[int]int
	LogProbs            *bool
	TopLogProbs         *int
	MaxCompletionTokens *int
	CompletionChoices   *int
	PresencePenalty     *float64
	Seed                *int
	Stop                []string
	Temperature         *float64
	TopP                *float64
	User                string
	Ctx                 context.Context
	// Tools               []openaiTool
}

type CompletionUsage struct {
	PromptTokens            int
	CompletionTokens        int
	TotalTokens             int
	CompletionTokensDetails map[string]any
}

type CompletionChoice struct {
	Index        int
	Message      m.Message
	FinishReason string
}

type CompletionError struct {
	Message string
	Type    string
}

type CompletionResponse struct {
	Id         string
	Object     string
	Created    int
	Model      string
	Choices    []CompletionChoice
	Usage      CompletionUsage
	Error      CompletionError
	StatusCode int
}

func (or CompletionResponse) Messages() []string {
	if c := or.Choices; c == nil || len(c) == 0 {
		return []string{}
	}

	messages := []string{}
	for _, c := range or.Choices {
		// check if it's a streaming request
		if c.Message.Content != "" {
			messages = append(messages, c.Message.Content)
		}
	}

	return messages
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
