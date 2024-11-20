package gollum

import (
	"time"

	gem "github.com/azr4e1/gollum/gemini"
	m "github.com/azr4e1/gollum/message"
	ll "github.com/azr4e1/gollum/ollama"
	oai "github.com/azr4e1/gollum/openai"
)

func (cr CompletionRequest) ToOpenAI() oai.CompletionRequest {
	messages := []oai.Message{}
	if system := cr.System.Content; system != "" {
		systemMessage := oai.Message{Role: "system", Content: system}
		messages = append(messages, systemMessage)
	}
	for _, mess := range cr.Messages {
		messages = append(messages, oai.Message{Role: mess.Role, Content: mess.Content})
	}
	// keep it simple stupid
	completionChoice := 1
	request := oai.CompletionRequest{
		Model:               cr.Model,
		Messages:            messages,
		Stream:              cr.Stream,
		FreqPenalty:         cr.FreqPenalty,
		LogitBias:           cr.LogitBias,
		LogProbs:            cr.LogProbs,
		TopLogProbs:         cr.TopLogProbs,
		MaxCompletionTokens: cr.MaxCompletionTokens,
		CompletionChoices:   &completionChoice,
		PresencePenalty:     cr.PresencePenalty,
		Seed:                cr.Seed,
		Stop:                cr.Stop,
		Temperature:         cr.Temperature,
		TopP:                cr.TopP,
		User:                cr.User,
		Ctx:                 cr.Ctx,
		// Tools:               []openaiTool,
	}

	return request
}

func (cr CompletionRequest) ToGemini() gem.CompletionRequest {
	messDict := map[string]string{
		"assistant": "model",
		"system":    "system",
		"user":      "user",
	}
	messages := []gem.Message{}
	for _, mess := range cr.Messages {
		part := [](map[string]string){
			{"text": mess.Content},
		}
		messages = append(messages, gem.Message{Role: messDict[mess.Role], Part: part})
	}
	var system map[string](map[string]string)
	if systemMessage := cr.System.Content; systemMessage != "" {
		system = map[string](map[string]string){
			"parts": {"text": systemMessage},
		}
	}

	var config map[string]any
	configOptions := map[string]any{
		"stopSequences":   cr.Stop,
		"temperature":     cr.Temperature,
		"maxOutputTokens": cr.MaxCompletionTokens,
		"topP":            cr.TopP,
		"topK":            cr.TopK,
	}
	for opt, val := range configOptions {
		if val != nil {
			if config == nil {
				config = make(map[string]any)
			}

			config[opt] = val
		}
	}

	request := gem.CompletionRequest{
		Model:         cr.Model,
		Messages:      messages,
		SystemMessage: system,
		Stream:        cr.Stream,
		Config:        config,
		Ctx:           cr.Ctx,
		// Tools:               []openaiTool,
	}

	return request
}

func (cr CompletionRequest) ToOllama() ll.CompletionRequest {
	messages := []ll.Message{}
	if system := cr.System.Content; system != "" {
		systemMessage := ll.Message{Role: "system", Content: system}
		messages = append(messages, systemMessage)
	}
	for _, mess := range cr.Messages {
		messages = append(messages, ll.Message{Role: mess.Role, Content: mess.Content})
	}
	request := ll.CompletionRequest{
		Model:    cr.Model,
		Messages: messages,
		Stream:   cr.Stream,
		Ctx:      cr.Ctx,
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

func ResponseFromGemini(response gem.CompletionResponse) CompletionResponse {
	usage := CompletionUsage{
		PromptTokens:     response.Usage.PromptTokens,
		CompletionTokens: response.Usage.CompletionTokens,
		TotalTokens:      response.Usage.TotalTokens,
	}

	message := m.Message{}
	finishReason := false
	if len(response.Choices) != 0 {
		c := response.Choices[0]
		if content := c.Content.Part[0]["text"]; c.Content.Role == "user" {
			message = m.UserMessage(content)
		} else {
			message = m.AssistantMessage(content)
		}

		if c.FinishReason != "" {
			finishReason = true
		}

	}

	var compErr CompletionError
	if response.Error.Status != "" {
		compErr = CompletionError{
			Message: response.Error.Message,
			Type:    response.Error.Status,
		}
	}
	converted := CompletionResponse{
		Model:      response.Model,
		Message:    message,
		Done:       finishReason,
		Usage:      usage,
		Error:      compErr,
		StatusCode: response.StatusCode,
	}

	return converted
}

func ResponseFromOpenAI(response oai.CompletionResponse) CompletionResponse {
	usage := CompletionUsage{
		PromptTokens:            response.Usage.PromptTokens,
		CompletionTokens:        response.Usage.CompletionTokens,
		TotalTokens:             response.Usage.TotalTokens,
		CompletionTokensDetails: response.Usage.CompletionTokensDetails,
	}

	message := m.Message{}
	finishReason := false
	if len(response.Choices) != 0 {
		c := response.Choices[0]
		role := ""
		content := ""
		if c.Message.Content != "" {
			role = c.Message.Role
			content = c.Message.Content
		} else if c.Delta.Content != "" {
			role = c.Delta.Role
			content = c.Delta.Content
		}

		if content != "" {
			if role == "user" {
				message = m.UserMessage(content)
			} else {
				message = m.AssistantMessage(content)
			}

		}

		if c.FinishReason != "" {
			finishReason = true
		}
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
		Message:    message,
		Done:       finishReason,
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

	message := m.Message{}
	var finishReason bool
	if content := response.Message.Content; content != "" {
		if response.Message.Role == "user" {
			message = m.UserMessage(content)
		} else {
			message = m.AssistantMessage(content)
		}
	}
	if response.Done {
		finishReason = true
	}

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
		Message:    message,
		Done:       finishReason,
		Usage:      usage,
		Error:      compErr,
		StatusCode: response.StatusCode,
	}

	return converted
}

func openaiComplete(request *CompletionRequest, c LLMClient) (CompletionRequest, CompletionResponse, error) {
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
	_, result, err := openaiClient.Complete(&openaiReq)
	if err != nil {
		return *request, CompletionResponse{}, err
	}

	return *request, ResponseFromOpenAI(result), nil
}

func geminiComplete(request *CompletionRequest, c LLMClient) (CompletionRequest, CompletionResponse, error) {
	geminiReq := request.ToGemini()
	geminiClient, err := c.ToGemini()
	if err != nil {
		return *request, CompletionResponse{}, err
	}
	if c.stream {
		streamFunc := func(geminiRes gem.CompletionResponse) error {
			res := ResponseFromGemini(geminiRes)
			return c.streamFunction(res)
		}
		geminiClient.EnableStream(streamFunc)
	}
	_, result, err := geminiClient.Complete(&geminiReq)
	if err != nil {
		return *request, CompletionResponse{}, err
	}

	return *request, ResponseFromGemini(result), nil
}

func ollamaComplete(request *CompletionRequest, c LLMClient) (CompletionRequest, CompletionResponse, error) {
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
	_, result, err := ollamaClient.Complete(&ollamaReq)
	if err != nil {
		return *request, CompletionResponse{}, err
	}

	return *request, ResponseFromOllama(result), nil
}
