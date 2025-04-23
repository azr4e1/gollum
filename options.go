package gollum

import (
	"context"
	"errors"
	m "github.com/azr4e1/gollum/message"
)

type clientOption func(*LLMClient) error
type completionOption func(*CompletionRequest) error
type speechOption func(*TTSRequest) error

func WithProvider(provider llmProvider) clientOption {
	return func(lc *LLMClient) error {
		lc.provider = provider

		return nil
	}
}
func WithAPIKey(apiKey string) clientOption {
	return func(lc *LLMClient) error {
		if apiKey == "" {
			return errors.New("must provide API key.")
		}
		lc.apiKey = apiKey

		return nil
	}
}
func WithAPIBase(apiBase string) clientOption {
	return func(lc *LLMClient) error {
		lc.apiBase = apiBase

		return nil
	}
}

func WithModel(modelName string) completionOption {
	return func(oR *CompletionRequest) error {
		oR.Model = modelName

		return nil
	}
}

func WithChat(chat m.Chat) completionOption {
	return func(oR *CompletionRequest) error {
		if chat.IsEmpty() {
			return errors.New("Missing messages to send.")
		}
		oR.System = chat.SystemMessage()
		oR.Messages = chat.History()

		return nil
	}
}

func WithMessage(message string) completionOption {
	return func(oR *CompletionRequest) error {
		if len(message) == 0 {
			return errors.New("Missing message to send.")
		}
		messages := []m.Message{m.UserMessage(message)}
		oR.Messages = messages

		return nil
	}
}

func WithFreqPenalty(freqPenalty float64) completionOption {
	return func(oR *CompletionRequest) error {
		if freqPenalty < -2.0 || freqPenalty > 2.0 {
			return errors.New("frequency penalty must be between -2.0 and 2.0.")
		}
		oR.FreqPenalty = &freqPenalty

		return nil
	}
}

func WithLogitBias(logitBias map[int]int) completionOption {
	return func(oR *CompletionRequest) error {
		if len(logitBias) == 0 {
			return errors.New("Map cannot be empty.")
		}
		oR.LogitBias = logitBias

		return nil
	}
}

func WithLogProbs(logProbs bool) completionOption {
	return func(oR *CompletionRequest) error {
		oR.LogProbs = &logProbs

		return nil
	}
}

func WithTopLogProbs(topLogProbs int) completionOption {
	return func(oR *CompletionRequest) error {
		if topLogProbs < 0 || topLogProbs > 20 {
			return errors.New("top_logprobs must be between 0 and 20.")
		}
		oR.TopLogProbs = &topLogProbs

		return nil
	}
}

func WithMaxCompletionTokens(maxCompletionTokens int) completionOption {
	return func(oR *CompletionRequest) error {
		if maxCompletionTokens <= 0 {
			return errors.New("Max completion tokens cannot be negative or zero.")
		}
		oR.MaxCompletionTokens = &maxCompletionTokens

		return nil
	}
}

func WithPresencePenalty(presencePenalty float64) completionOption {
	return func(oR *CompletionRequest) error {
		if presencePenalty < -2.0 || presencePenalty > 2.0 {
			return errors.New("Presence penalty must be between -2.0 and 2.0")
		}
		oR.PresencePenalty = &presencePenalty

		return nil
	}
}

func WithSeed(seed int) completionOption {
	return func(oR *CompletionRequest) error {
		oR.Seed = &seed

		return nil
	}
}

func WithStop(stop []string) completionOption {
	return func(oR *CompletionRequest) error {
		oR.Stop = stop

		return nil
	}
}

func WithTemperature(temperature float64) completionOption {
	return func(oR *CompletionRequest) error {
		if temperature < 0.0 || temperature > 2.0 {
			return errors.New("Temperature must be between 0.0 and 2.0.")
		}
		oR.Temperature = &temperature

		return nil
	}
}

func WithTopP(topP float64) completionOption {
	return func(oR *CompletionRequest) error {
		if topP < 0.0 || topP > 1 {
			return errors.New("Top P must be between 0 and 1.")
		}
		oR.TopP = &topP

		return nil
	}
}

func WithTopK(topK int) completionOption {
	return func(oR *CompletionRequest) error {
		if topK < 0 {
			return errors.New("Top K must be between greater than 0.")
		}
		oR.TopK = &topK

		return nil
	}
}

func WithUser(user string) completionOption {
	return func(oR *CompletionRequest) error {
		if user == "" {
			return errors.New("Cannot set user to empty string.")
		}
		oR.User = user

		return nil
	}
}

func WithContext(ctx context.Context) completionOption {
	return func(oR *CompletionRequest) error {
		oR.Ctx = ctx

		return nil
	}
}

func WithTool(tools ...Tool) completionOption {
	return func(oR *CompletionRequest) error {
		oR.Tools = tools

		return nil
	}
}

func WithTTSModel(model string) speechOption {
	return func(aR *TTSRequest) error {
		if model == "" {
			return errors.New("model is missing")
		}
		aR.Model = model
		return nil
	}
}
func WithTTSInput(input string) speechOption {
	return func(aR *TTSRequest) error {
		if input == "" {
			return errors.New("input is empty.")
		}
		aR.Input = input
		return nil
	}
}
func WithTTSVoice(voice string) speechOption {
	return func(aR *TTSRequest) error {
		if voice == "" {
			return errors.New("voice is missing")
		}
		aR.Voice = voice
		return nil
	}
}
func WithTTSFormat(format string) speechOption {
	return func(aR *TTSRequest) error {
		if format == "" {
			return errors.New("format is missing")
		}
		aR.Format = format
		return nil
	}
}
func WithTTSSpeed(speed float64) speechOption {
	return func(aR *TTSRequest) error {
		if speed < 0.25 || speed > 4 {
			return errors.New("speed must be between 0.25 and 4. Default is 1.")
		}
		aR.Speed = &speed
		return nil
	}
}

func WithTTSContext(ctx context.Context) speechOption {
	return func(aR *TTSRequest) error {
		aR.Ctx = ctx

		return nil
	}
}
