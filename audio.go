package gollum

import "errors"

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

func WithAudioModel(model string) speechOption {
	return func(aR *AudioRequest) error {
		if model == "" {
			return errors.New("model is missing")
		}
		aR.Model = model
		return nil
	}
}
func WithAudioInput(input string) speechOption {
	return func(aR *AudioRequest) error {
		if input == "" {
			return errors.New("input is empty.")
		}
		aR.Input = input
		return nil
	}
}
func WithAudioVoice(voice string) speechOption {
	return func(aR *AudioRequest) error {
		if voice == "" {
			return errors.New("voice is missing")
		}
		aR.Voice = voice
		return nil
	}
}
func WithAudioFormat(format string) speechOption {
	return func(aR *AudioRequest) error {
		if format == "" {
			return errors.New("format is missing")
		}
		aR.Format = format
		return nil
	}
}
func WithAudioSpeed(speed float64) speechOption {
	return func(aR *AudioRequest) error {
		if speed < 0.25 || speed > 4 {
			return errors.New("speed must be between 0.25 and 4. Default is 1.")
		}
		aR.Speed = &speed
		return nil
	}
}

func NewAudioRequest(opts ...speechOption) (*AudioRequest, error) {
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
