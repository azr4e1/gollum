package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	m "github.com/azr4e1/gollum/message"
	"net/http"
)

type CompletionRequest struct {
	Model               string          `json:"model"`
	Messages            []m.Message     `json:"messages"`
	Stream              bool            `json:"stream"`
	Tools               []openaiTool    `json:"tools,omitempty"`
	FreqPenalty         *float64        `json:"frequency_penalty,omitempty"`
	LogitBias           map[int]int     `json:"logit_bias,omitempty"`
	LogProbs            *bool           `json:"logprobs,omitempty"`
	TopLogProbs         *int            `json:"top_logprobs,omitempty"`
	MaxCompletionTokens *int            `json:"max_completion_tokens,omitempty"`
	CompletionChoices   *int            `json:"n,omitempty"`
	PresencePenalty     *float64        `json:"presence_penalty,omitempty"`
	Seed                *int            `json:"seed,omitempty"`
	Stop                []string        `json:"stop,omitempty"`
	Temperature         *float64        `json:"temperature,omitempty"`
	TopP                *float64        `json:"top_p,omitempty"`
	User                string          `json:"user,omitempty"`
	Ctx                 context.Context `json:"-"`
}

type CompletionUsage struct {
	PromptTokens            int            `json:"prompt_tokens"`
	CompletionTokens        int            `json:"completion_tokens"`
	TotalTokens             int            `json:"total_tokens"`
	CompletionTokensDetails map[string]any `json:"completion_tokens_details"`
}

type CompletionChoice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message,omitempty"`
	Delta        Message `json:"delta,omitempty"`
	FinishReason string  `json:"finish_reason"`
}

type CompletionError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type CompletionResponse struct {
	Id         string             `json:"id"`
	Object     string             `json:"object"`
	Created    int                `json:"created"`
	Model      string             `json:"model"`
	Choices    []CompletionChoice `json:"choices"`
	Usage      CompletionUsage    `json:"usage"`
	Error      CompletionError    `json:"error,omitempty"`
	StatusCode int                `json:"status_code"`
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
		} else if c.Delta.Content != "" {
			messages = append(messages, c.Delta.Content)
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

func makeHTTPCompletionRequest(request *CompletionRequest, oc OpenaiClient) (*http.Response, error) {
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, completionURL, bytes.NewReader(jsonRequest))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", oc.apiKey))
	if request.Ctx != nil {
		req = req.WithContext(request.Ctx)
	}

	client := http.Client{Timeout: oc.Timeout}
	res, err := client.Do(req)

	return res, err
}
