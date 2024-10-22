package gollum

import "errors"

const streamEnd = "gollum_end_of_stream"

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

func (or CompletionResponse) GetMessages() []string {
	if c := or.Choices; c == nil || len(c) == 0 {
		return []string{}
	}

	messages := []string{}
	for _, c := range or.Choices {
		// check if it's a streaming request
		if c.Message != nil && c.Message.Content != "" {
			messages = append(messages, c.Message.Content)
		} else if c.Delta != nil && c.Delta.Content != "" {
			messages = append(messages, c.Delta.Content)
		}
	}

	return messages
}

func newEOS(message string) CompletionResponse {
	return CompletionResponse{Error: &CompletionError{Type: streamEnd, Message: message}}
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
