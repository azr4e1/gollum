package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	m "github.com/azr4e1/gollum/message"
	"net/http"
	"net/url"
	"path"
)

type CompletionRequest struct {
	Model    string          `json:"model"`
	Messages []m.Message     `json:"messages"`
	Stream   bool            `json:"stream"`
	Ctx      context.Context `json:"-"`
	// Tools               []openaiTool `json:"tools,omitempty"`
	// FreqPenalty         *float64     `json:"frequency_penalty,omitempty"`
	// LogitBias           map[int]int  `json:"logit_bias,omitempty"`
	// LogProbs            *bool        `json:"logprobs,omitempty"`
	// TopLogProbs         *int         `json:"top_logprobs,omitempty"`
	// MaxCompletionTokens *int         `json:"max_completion_tokens,omitempty"`
	// CompletionChoices   *int         `json:"n,omitempty"`
	// PresencePenalty     *float64     `json:"presence_penalty,omitempty"`
	// Seed                *int         `json:"seed,omitempty"`
	// Stop                []string     `json:"stop,omitempty"`
	// Temperature         *float64     `json:"temperature,omitempty"`
	// TopP                *float64     `json:"top_p,omitempty"`
	// User                string       `json:"user,omitempty"`
}

type CompletionResponse struct {
	Created            string    `json:"created_at"`
	Model              string    `json:"model"`
	Message            m.Message `json:"message"`
	Done               bool      `json:"done"`
	TotalDuration      int       `json:"total_duration"`
	LoadDuration       int       `json:"load_duration"`
	PromptEvalCount    int       `json:"prompt_eval_count"`
	PromptEvalDuration int       `json:"prompt_eval_duration"`
	EvalCount          int       `json:"eval_count"`
	EvalDuration       int       `json:"eval_duration"`
	Error              string    `json:"error,omitempty"`
	StatusCode         int       `json:"status_code"`
}

func (or CompletionResponse) err() error {
	if or.Error == "" {
		return nil
	}
	return errors.New(or.Error)
}

func makeHTTPCompletionRequest(request *CompletionRequest, oc OllamaClient) (*http.Response, error) {
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	url, err := url.Parse(oc.baseURL)
	if err != nil {
		return nil, err
	}
	url.Path = path.Join(url.Path, completionURL)
	req, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewReader(jsonRequest))
	if err != nil {
		return nil, err
	}
	if request.Ctx != nil {
		req = req.WithContext(request.Ctx)
	}

	client := http.Client{Timeout: oc.Timeout}
	res, err := client.Do(req)

	return res, err
}
