package gollum

import (
	"context"
	"errors"
	"fmt"
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

func NewTTSRequest(opts ...speechOption) (*TTSRequest, error) {
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
