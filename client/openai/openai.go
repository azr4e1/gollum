package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	BASEURL = "https://api.openai.com/v1/chat/completions"
)

type openaiRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type openaiResponse struct {
	Id         string         `json:"id"`
	Object     string         `json:"object"`
	Created    int            `json:"created"`
	Model      string         `json:"model"`
	Choices    []openaiChoice `json:"choices"`
	Usage      openaiUsage    `json:"usage"`
	Error      openaiError    `json:"error"`
	StatusCode int            `json:"status_code"`
}

type openaiUsage struct {
	PromptTokens            int            `json:"prompt_tokens"`
	CompletionTokens        int            `json:"completion_tokens"`
	TotalTokens             int            `json:"total_tokens"`
	CompletionTokensDetails map[string]any `json:"completion_tokens_details"`
}

type openaiChoice struct {
	Index   int     `json:"index"`
	Message message `json:"message"`
}

type openaiError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type openaiClient struct {
	apiKey string
}

func NewOpenaiClient(apiKey string) openaiClient {
	return openaiClient{apiKey: apiKey}
}

func (oc openaiClient) Complete(model string, chat llmChat) (openaiRequest, openaiResponse, error) {

	request := openaiRequest{
		Model:    model,
		Messages: chat.GetHistory(),
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return openaiRequest{}, openaiResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, BASEURL, bytes.NewReader(jsonRequest))
	if err != nil {
		return openaiRequest{}, openaiResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", oc.apiKey))

	client := http.Client{Timeout: time.Duration(30 * time.Second)}
	res, err := client.Do(req)

	if err != nil {
		return openaiRequest{}, openaiResponse{}, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return openaiRequest{}, openaiResponse{}, err
	}

	openaiRes := new(openaiResponse)
	json.Unmarshal(body, openaiRes)

	// attach status code to response object
	openaiRes.StatusCode = res.StatusCode

	return request, *openaiRes, nil
}
