package gollum

import (
	oai "github.com/azr4e1/gollum/openai"
)

func (ar TTSRequest) ToOpenAI() oai.TTSRequest {
	ttsReq := oai.TTSRequest{
		Model:  ar.Model,
		Input:  ar.Input,
		Voice:  ar.Voice,
		Format: ar.Format,
		Speed:  ar.Speed,
		Ctx:    ar.Ctx,
	}

	return ttsReq
}

func SpeechResponseFromOpenAI(response oai.TTSResponse) TTSResponse {
	var error TTSError
	if response.Err() != nil {
		error = TTSError{
			Message: response.Error.Message,
			Type:    response.Error.Type,
		}
	}
	ttsResponse := TTSResponse{
		Audio:      response.Audio,
		Error:      error,
		StatusCode: response.StatusCode,
	}

	return ttsResponse
}

func openaiTTS(request *TTSRequest, c LLMClient) (TTSRequest, TTSResponse, error) {
	openaiReq := request.ToOpenAI()
	openaiClient, err := c.ToOpenAI()
	if err != nil {
		return *request, TTSResponse{}, err
	}
	_, result, err := openaiClient.TextToSpeech(&openaiReq)
	if err != nil {
		return *request, TTSResponse{}, err
	}

	return *request, SpeechResponseFromOpenAI(result), nil
}
