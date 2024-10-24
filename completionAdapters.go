package gollum

import (
	"time"

	ll "github.com/azr4e1/gollum/ollama"
	oai "github.com/azr4e1/gollum/openai"
)

func (cr CompletionRequest) ToOpenAI() oai.CompletionRequest {
	messages := []oai.Message{}
	for _, m := range cr.Messages {
		messages = append(messages, oai.Message{Role: m.Role, Content: m.Content})
	}
	request := oai.CompletionRequest{
		Model:               cr.Model,
		Messages:            messages,
		Stream:              cr.Stream,
		FreqPenalty:         cr.FreqPenalty,
		LogitBias:           cr.LogitBias,
		LogProbs:            cr.LogProbs,
		TopLogProbs:         cr.TopLogProbs,
		MaxCompletionTokens: cr.MaxCompletionTokens,
		CompletionChoices:   cr.CompletionChoices,
		PresencePenalty:     cr.PresencePenalty,
		Seed:                cr.Seed,
		Stop:                cr.Stop,
		Temperature:         cr.Temperature,
		TopP:                cr.TopP,
		User:                cr.User,
		// Tools:               []openaiTool,
	}

	return request
}

func (cr CompletionRequest) ToOllama() ll.CompletionRequest {
	messages := []ll.Message{}
	for _, m := range cr.Messages {
		messages = append(messages, ll.Message{Role: m.Role, Content: m.Content})
	}
	request := ll.CompletionRequest{
		Model:    cr.Model,
		Messages: messages,
		Stream:   cr.Stream,
		// FreqPenalty:         cr.FreqPenalty,
		// LogitBias:           cr.LogitBias,
		// LogProbs:            cr.LogProbs,
		// TopLogProbs:         cr.TopLogProbs,
		// MaxCompletionTokens: cr.MaxCompletionTokens,
		// CompletionChoices:   cr.CompletionChoices,
		// PresencePenalty:     cr.PresencePenalty,
		// Seed:                cr.Seed,
		// Stop:                cr.Stop,
		// Temperature:         cr.Temperature,
		// TopP:                cr.TopP,
		// User:                cr.User,
		// Tools:               []openaiTool,
	}

	return request
}

func ResponseFromOpenAI(response oai.CompletionResponse) CompletionResponse {
	usage := CompletionUsage{
		PromptTokens:            response.Usage.PromptTokens,
		CompletionTokens:        response.Usage.CompletionTokens,
		TotalTokens:             response.Usage.TotalTokens,
		CompletionTokensDetails: response.Usage.CompletionTokensDetails,
	}

	choices := []CompletionChoice{}
	for _, c := range response.Choices {
		var message Message
		if c.Message.Content != "" {
			message = Message{Role: c.Message.Role, Content: c.Message.Content}
		} else if c.Delta.Content != "" {
			message = Message{Role: c.Delta.Role, Content: c.Delta.Content}
		}
		choice := CompletionChoice{
			Index:        c.Index,
			Message:      message,
			FinishReason: c.FinishReason,
		}
		choices = append(choices, choice)
	}

	var compErr CompletionError
	if response.Error.Type != "" {
		compErr = CompletionError{
			Message: response.Error.Message,
			Type:    response.Error.Type,
		}
	}
	converted := CompletionResponse{
		Id:         response.Id,
		Object:     response.Object,
		Created:    response.Created,
		Model:      response.Model,
		Choices:    choices,
		Usage:      usage,
		Error:      compErr,
		StatusCode: response.StatusCode,
	}

	return converted
}

func ResponseFromOllama(response ll.CompletionResponse) CompletionResponse {
	usage := CompletionUsage{
		PromptTokens:     response.PromptEvalCount,
		CompletionTokens: response.EvalCount,
		TotalTokens:      response.PromptEvalCount + response.EvalCount,
	}

	choices := []CompletionChoice{}
	var message Message
	if response.Message.Content != "" {
		message = Message{Role: response.Message.Role, Content: response.Message.Content}
	}
	var finishReason string
	if response.Done {
		finishReason = "done"
	}
	choice := CompletionChoice{
		Message:      message,
		FinishReason: finishReason,
	}
	choices = append(choices, choice)

	var compErr CompletionError
	if response.Error != "" {
		compErr = CompletionError{
			Message: response.Error,
		}
	}

	var created time.Time
	if response.Created != "" {
		var timeErr error
		created, timeErr = time.Parse(time.RFC3339Nano, response.Created)
		if timeErr != nil {
			panic(timeErr)
		}
	}
	converted := CompletionResponse{
		Created:    int(created.Unix()),
		Model:      response.Model,
		Choices:    choices,
		Usage:      usage,
		Error:      compErr,
		StatusCode: response.StatusCode,
	}

	return converted
}

func openaiComplete(request *CompletionRequest, c llmClient) (CompletionRequest, CompletionResponse, error) {
	openaiReq := request.ToOpenAI()
	openaiClient, err := c.ToOpenAI()
	if err != nil {
		return *request, CompletionResponse{}, err
	}
	if c.stream {
		streamFunc := func(openaiRes oai.CompletionResponse) error {
			res := ResponseFromOpenAI(openaiRes)
			return c.streamFunction(res)
		}
		openaiClient.EnableStream(streamFunc)
	}
	_, result, err := openaiClient.CompleteWithCustomRequest(&openaiReq)
	if err != nil {
		return *request, CompletionResponse{}, err
	}

	return *request, ResponseFromOpenAI(result), nil
}

func ollamaComplete(request *CompletionRequest, c llmClient) (CompletionRequest, CompletionResponse, error) {
	ollamaReq := request.ToOllama()
	ollamaClient, err := c.ToOllama()
	if err != nil {
		return *request, CompletionResponse{}, err
	}
	if c.stream {
		streamFunc := func(ollamaRes ll.CompletionResponse) error {
			res := ResponseFromOllama(ollamaRes)
			return c.streamFunction(res)
		}
		ollamaClient.EnableStream(streamFunc)
	}
	_, result, err := ollamaClient.CompleteWithCustomRequest(&ollamaReq)
	if err != nil {
		return *request, CompletionResponse{}, err
	}

	return *request, ResponseFromOllama(result), nil
}
