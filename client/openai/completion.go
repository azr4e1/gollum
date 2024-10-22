package openai

type completionRequest struct {
	Model               string       `json:"model"`
	Messages            []message    `json:"messages"`
	Stream              bool         `json:"stream"`
	Tools               []openaiTool `json:"tools,omitempty"`
	FreqPenalty         *float64     `json:"frequency_penalty,omitempty"`
	LogitBias           map[int]int  `json:"logit_bias,omitempty"`
	LogProbs            *bool        `json:"logprobs,omitempty"`
	TopLogProbs         *int         `json:"top_logprobs,omitempty"`
	MaxCompletionTokens *int         `json:"max_completion_tokens,omitempty"`
	CompletionChoices   *int         `json:"n,omitempty"`
	PresencePenalty     *float64     `json:"presence_penalty,omitempty"`
	Seed                *int         `json:"seed,omitempty"`
	Stop                []string     `json:"stop,omitempty"`
	Temperature         *float64     `json:"temperature,omitempty"`
	TopP                *float64     `json:"top_p,omitempty"`
	User                string       `json:"user,omitempty"`
}

type completionUsage struct {
	PromptTokens            int            `json:"prompt_tokens"`
	CompletionTokens        int            `json:"completion_tokens"`
	TotalTokens             int            `json:"total_tokens"`
	CompletionTokensDetails map[string]any `json:"completion_tokens_details"`
}

type completionChoice struct {
	Index        int      `json:"index"`
	Message      *message `json:"message,omitempty"`
	Delta        *message `json:"delta,omitempty"`
	FinishReason string   `json:"finish_reason"`
}

type completionError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type completionResponse struct {
	Id         string             `json:"id"`
	Object     string             `json:"object"`
	Created    int                `json:"created"`
	Model      string             `json:"model"`
	Choices    []completionChoice `json:"choices"`
	Usage      completionUsage    `json:"usage"`
	Error      *completionError   `json:"error,omitempty"`
	StatusCode int                `json:"status_code"`
}

func (or completionResponse) GetMessages() []string {
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

func (or completionResponse) IsEOS() bool {
	return or.Error.Type == streamEnd
}

func newEOS(message string) completionResponse {
	return completionResponse{Error: &completionError{Type: streamEnd, Message: message}}
}
