package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
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
	Model  string   `json:"model"`
	Input  string   `json:"input"`
	Voice  string   `json:"voice"`
	Format string   `json:"response_format,omitempty"`
	Speed  *float64 `json:"speed,omitempty"`
}

type TTSResponse struct {
	Audio      []byte    `json:"audio"`
	Error      *TTSError `json:"error,omitempty"`
	StatusCode int       `json:"status_code"`
}

type TTSError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func NewTTSRequest(opts ...ttsOption) (*TTSRequest, error) {
	request := new(TTSRequest)
	for _, o := range opts {
		err := o(request)
		if err != nil {
			return &TTSRequest{}, err
		}
	}

	if request.Model == "" {
		return &TTSRequest{}, errors.New("missing model.")
	}
	if request.Voice == "" {
		return &TTSRequest{}, errors.New("missing voice.")
	}
	if request.Input == "" {
		return &TTSRequest{}, errors.New("missing input.")
	}

	return request, nil
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

	client := http.Client{Timeout: time.Duration(30 * time.Second)}
	res, err := client.Do(req)

	return res, err
}
