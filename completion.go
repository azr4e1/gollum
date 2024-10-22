package gollum

import "errors"

type CompletionRequest struct {
	Model               string
	Messages            []Message
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
	Message      *Message
	Delta        *Message
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
	Error      *CompletionError
	StatusCode int
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
