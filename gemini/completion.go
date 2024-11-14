package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	// m "github.com/azr4e1/gollum/message"
	"net/http"
)

const (
	nonStream = "generateContent"
	stream    = "streamGenerateContent"
)

type CompletionRequest struct {
	Model         string                         `json:"-"`
	Messages      []Message                      `json:"contents"`
	SystemMessage map[string](map[string]string) `json:"system_instruction,omitempty"`
	Stream        bool                           `json:"-"`
	// Tools               []openaiTool    `json:"tools,omitempty"`
	// FreqPenalty         *float64        `json:"frequency_penalty,omitempty"`
	// LogitBias           map[int]int     `json:"logit_bias,omitempty"`
	// LogProbs            *bool           `json:"logprobs,omitempty"`
	// TopLogProbs         *int            `json:"top_logprobs,omitempty"`
	// MaxCompletionTokens *int            `json:"max_completion_tokens,omitempty"`
	// CompletionChoices   *int            `json:"n,omitempty"`
	// PresencePenalty     *float64        `json:"presence_penalty,omitempty"`
	// Seed                *int            `json:"seed,omitempty"`
	// Stop                []string        `json:"stop,omitempty"`
	// Temperature         *float64        `json:"temperature,omitempty"`
	// TopP                *float64        `json:"top_p,omitempty"`
	// User                string          `json:"user,omitempty"`
	Ctx context.Context `json:"-"`
}

type CompletionUsage struct {
	PromptTokens     int `json:"promptTokenCount"`
	CompletionTokens int `json:"candidatesTokenCount"`
	TotalTokens      int `json:"totalTokenCount"`
}

type SafetyRating struct {
	Category    string `json:"category"`
	Probability string `json:"probability"`
}

type CompletionChoice struct {
	Content       Message        `json:"content,omitempty"`
	FinishReason  string         `json:"finishReason"`
	Index         int            `json:"index"`
	SafetyRatings []SafetyRating `json:"safetyRatings"`
}

type CompletionError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

type CompletionResponse struct {
	Model      string             `json:"modelVersion"`
	Choices    []CompletionChoice `json:"candidates"`
	Usage      CompletionUsage    `json:"usageMetadata"`
	Error      CompletionError    `json:"error,omitempty"`
	StatusCode int                `json:"status_code"`
}

func (or CompletionResponse) err() error {
	if or.Error.Status == "" && or.Error.Message == "" {
		return nil
	}
	return errors.New(fmt.Sprintf("%s: %s", or.Error.Status, or.Error.Message))
}

func makeHTTPCompletionRequest(request *CompletionRequest, oc GeminiClient) (*http.Response, error) {
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	streamQuery := nonStream
	if oc.stream {
		streamQuery = stream
	}
	fullURL := fmt.Sprintf(completionURL, request.Model, streamQuery, oc.apiKey)
	req, err := http.NewRequest(http.MethodPost, fullURL, bytes.NewReader(jsonRequest))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if request.Ctx != nil {
		req = req.WithContext(request.Ctx)
	}

	client := http.Client{Timeout: oc.Timeout}
	res, err := client.Do(req)

	return res, err
}
