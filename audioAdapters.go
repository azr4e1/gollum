package gollum

import (
	oai "github.com/azr4e1/gollum/openai"
)

func (ar AudioRequest) ToOpenAI() oai.AudioRequest {
	audioReq := oai.AudioRequest{
		Model:  ar.Model,
		Input:  ar.Input,
		Voice:  ar.Voice,
		Format: ar.Format,
		Speed:  ar.Speed,
	}

	return audioReq
}

func SpeechResponseFromOpenAI(response oai.AudioResponse) AudioResponse {
	var error *AudioError
	if response.Error != nil {
		error = &AudioError{
			Message: response.Error.Message,
			Type:    response.Error.Type,
		}
	}
	audioResponse := AudioResponse{
		Audio:      response.Audio,
		Error:      error,
		StatusCode: response.StatusCode,
	}

	return audioResponse
}

func openaiSpeech(request *AudioRequest, c llmClient) (AudioRequest, AudioResponse, error) {
	openaiReq := request.ToOpenAI()
	openaiClient, err := c.ToOpenAI()
	if err != nil {
		return *request, AudioResponse{}, err
	}
	_, result, err := openaiClient.SpeechWithCustomRequest(&openaiReq)
	if err != nil {
		return *request, AudioResponse{}, err
	}

	return *request, SpeechResponseFromOpenAI(result), nil
}
