package gollum

import (
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

func ResponseFromOpenAI(response oai.CompletionResponse) CompletionResponse {
	usage := CompletionUsage{
		PromptTokens:            response.Usage.PromptTokens,
		CompletionTokens:        response.Usage.CompletionTokens,
		TotalTokens:             response.Usage.TotalTokens,
		CompletionTokensDetails: response.Usage.CompletionTokensDetails,
	}

	choices := []CompletionChoice{}
	for _, c := range response.Choices {
		var message *Message
		var delta *Message
		if c.Message != nil {
			message = &Message{Role: c.Message.Role, Content: c.Message.Content}
		}
		if c.Delta != nil {
			delta = &Message{Role: c.Delta.Role, Content: c.Delta.Content}
		}
		choice := CompletionChoice{
			Index:        c.Index,
			Message:      message,
			Delta:        delta,
			FinishReason: c.FinishReason,
		}
		choices = append(choices, choice)
	}

	var error *CompletionError
	if response.Error != nil {
		error = &CompletionError{
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
		Error:      error,
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
