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
type audioOption func(*audioRequest) error

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

type audioRequest struct {
	Model  audioModel  `json:"model"`
	Input  string      `json:"input"`
	Voice  openaiVoice `json:"voice"`
	Format audioFormat `json:"response_format,omitempty"`
	Speed  *float64    `json:"speed,omitempty"`
}

type audioResponse struct {
	Audio      []byte      `json:"audio"`
	Error      *audioError `json:"error,omitempty"`
	StatusCode int         `json:"status_code"`
}

type audioError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func WithAudioModel(model audioModel) audioOption {
	return func(aR *audioRequest) error {
		aR.Model = model
		return nil
	}
}
func WithAudioInput(input string) audioOption {
	return func(aR *audioRequest) error {
		if input == "" {
			return errors.New("input is empty.")
		}
		aR.Input = input
		return nil
	}
}
func WithAudioVoice(voice openaiVoice) audioOption {
	return func(aR *audioRequest) error {
		aR.Voice = voice
		return nil
	}
}
func WithAudioFormat(format audioFormat) audioOption {
	return func(aR *audioRequest) error {
		aR.Format = format
		return nil
	}
}
func WithAudioSpeed(speed float64) audioOption {
	return func(aR *audioRequest) error {
		if speed < 0.25 || speed > 4 {
			return errors.New("speed must be between 0.25 and 4. Default is 1.")
		}
		aR.Speed = &speed
		return nil
	}
}

func NewAudioRequest(opts ...audioOption) (*audioRequest, error) {
	request := new(audioRequest)
	for _, o := range opts {
		err := o(request)
		if err != nil {
			return nil, err
		}
	}

	if request.Model == "" {
		return nil, errors.New("missing model.")
	}
	if request.Voice == "" {
		return nil, errors.New("missing voice.")
	}
	if request.Input == "" {
		return nil, errors.New("missing input.")
	}

	return request, nil
}

func makeHTTPAudioRequest(request *audioRequest, oc OpenaiClient) (*http.Response, error) {
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, audioURL, bytes.NewReader(jsonRequest))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", oc.apiKey))

	client := http.Client{Timeout: time.Duration(30 * time.Second)}
	res, err := client.Do(req)

	return res, err
}
