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
type audioFormat string
type audioModel string
type audioOption func(*AudioRequest) error

const (
	Alloy   openaiVoice = "alloy"
	Echo    openaiVoice = "echo"
	Fable   openaiVoice = "fable"
	Onyx    openaiVoice = "onyx"
	Nova    openaiVoice = "nova"
	Shimmer openaiVoice = "shimmer"
)

const (
	MP3  audioFormat = "mp3"
	OPUS audioFormat = "opus"
	AAC  audioFormat = "aac"
	FLAC audioFormat = "flac"
	WAV  audioFormat = "wav"
	PCM  audioFormat = "pcm"
)

const (
	TTS1   audioModel = "tts-1"
	TTS1HD audioModel = "tts-1-hd"
)

type AudioRequest struct {
	Model  string   `json:"model"`
	Input  string   `json:"input"`
	Voice  string   `json:"voice"`
	Format string   `json:"response_format,omitempty"`
	Speed  *float64 `json:"speed,omitempty"`
}

type AudioResponse struct {
	Audio      []byte      `json:"audio"`
	Error      *AudioError `json:"error,omitempty"`
	StatusCode int         `json:"status_code"`
}

type AudioError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func WithAudioModel(model audioModel) audioOption {
	return func(aR *AudioRequest) error {
		aR.Model = string(model)
		return nil
	}
}
func WithAudioInput(input string) audioOption {
	return func(aR *AudioRequest) error {
		if input == "" {
			return errors.New("input is empty.")
		}
		aR.Input = input
		return nil
	}
}
func WithAudioVoice(voice openaiVoice) audioOption {
	return func(aR *AudioRequest) error {
		aR.Voice = string(voice)
		return nil
	}
}
func WithAudioFormat(format audioFormat) audioOption {
	return func(aR *AudioRequest) error {
		aR.Format = string(format)
		return nil
	}
}
func WithAudioSpeed(speed float64) audioOption {
	return func(aR *AudioRequest) error {
		if speed < 0.25 || speed > 4 {
			return errors.New("speed must be between 0.25 and 4. Default is 1.")
		}
		aR.Speed = &speed
		return nil
	}
}

func NewAudioRequest(opts ...audioOption) (*AudioRequest, error) {
	request := new(AudioRequest)
	for _, o := range opts {
		err := o(request)
		if err != nil {
			return &AudioRequest{}, err
		}
	}

	if request.Model == "" {
		return &AudioRequest{}, errors.New("missing model.")
	}
	if request.Voice == "" {
		return &AudioRequest{}, errors.New("missing voice.")
	}
	if request.Input == "" {
		return &AudioRequest{}, errors.New("missing input.")
	}

	return request, nil
}

func makeHTTPAudioRequest(request *AudioRequest, oc OpenaiClient) (*http.Response, error) {
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
