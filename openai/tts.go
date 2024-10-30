package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type openaiVoice string
type ttsFormat string
type ttsModel string
type ttsOption func(*TTSRequest) error

const (
	Alloy   openaiVoice = "alloy"
	Echo    openaiVoice = "echo"
	Fable   openaiVoice = "fable"
	Onyx    openaiVoice = "onyx"
	Nova    openaiVoice = "nova"
	Shimmer openaiVoice = "shimmer"
)

const (
	MP3  ttsFormat = "mp3"
	OPUS ttsFormat = "opus"
	AAC  ttsFormat = "aac"
	FLAC ttsFormat = "flac"
	WAV  ttsFormat = "wav"
	PCM  ttsFormat = "pcm"
)

const (
	TTS1   ttsModel = "tts-1"
	TTS1HD ttsModel = "tts-1-hd"
)

type TTSRequest struct {
	Model  string          `json:"model"`
	Input  string          `json:"input"`
	Voice  string          `json:"voice"`
	Format string          `json:"response_format,omitempty"`
	Speed  *float64        `json:"speed,omitempty"`
	Ctx    context.Context `json:"-"`
}

type TTSResponse struct {
	Audio      []byte   `json:"audio"`
	Error      TTSError `json:"error,omitempty"`
	StatusCode int      `json:"status_code"`
}

func (ttsr TTSResponse) Err() error {
	if ttsr.Error.Type == "" && ttsr.Error.Message == "" {
		return nil
	}

	return errors.New(fmt.Sprintf("%s: %s", ttsr.Error.Type, ttsr.Error.Message))
}

type TTSError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func makeHTTPTTSRequest(request *TTSRequest, oc OpenaiClient) (*http.Response, error) {
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, speechURL, bytes.NewReader(jsonRequest))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", oc.apiKey))
	if request.Ctx != nil {
		req = req.WithContext(request.Ctx)
	}

	client := http.Client{Timeout: oc.Timeout}
	res, err := client.Do(req)

	return res, err
}
