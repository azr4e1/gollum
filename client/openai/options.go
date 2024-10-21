package openai

import "errors"

type completionOption func(*openaiRequest) error

func WithModel(modelName string) completionOption {
	return func(oR *openaiRequest) error {
		oR.Model = modelName

		return nil
	}
}

func WithMessages(messages []message) completionOption {
	return func(oR *openaiRequest) error {
		if messages == nil || len(messages) == 0 {
			return errors.New("Missing messages to send.")
		}
		oR.Messages = messages

		return nil
	}
}

func WithFreqPenalty(freqPenalty float64) completionOption {
	return func(oR *openaiRequest) error {
		if freqPenalty < -2.0 || freqPenalty > 2.0 {
			return errors.New("frequency penalty must be between -2.0 and 2.0.")
		}
		oR.FreqPenalty = &freqPenalty

		return nil
	}
}

func WithLogitBias(logitBias map[int]int) completionOption {
	return func(oR *openaiRequest) error {
		if len(logitBias) == 0 {
			return errors.New("Map cannot be empty.")
		}
		oR.LogitBias = logitBias

		return nil
	}
}

func WithLogProbs(logProbs bool) completionOption {
	return func(oR *openaiRequest) error {
		oR.LogProbs = &logProbs

		return nil
	}
}

func WithTopLogProbs(topLogProbs int) completionOption {
	return func(oR *openaiRequest) error {
		if topLogProbs < 0 || topLogProbs > 20 {
			return errors.New("top_logprobs must be between 0 and 20.")
		}
		oR.TopLogProbs = &topLogProbs

		return nil
	}
}

func WithMaxCompletionTokens(maxCompletionTokens int) completionOption {
	return func(oR *openaiRequest) error {
		if maxCompletionTokens <= 0 {
			return errors.New("Max completion tokens cannot be negative or zero.")
		}
		oR.MaxCompletionTokens = &maxCompletionTokens

		return nil
	}
}

func WithCompletionChoices(completionChoices int) completionOption {
	return func(oR *openaiRequest) error {
		if completionChoices <= 0 {
			return errors.New("Number of completion choices cannot be negative or zero.")
		}
		oR.CompletionChoices = &completionChoices

		return nil
	}
}

func WithPresencePenalty(presencePenalty float64) completionOption {
	return func(oR *openaiRequest) error {
		if presencePenalty < -2.0 || presencePenalty > 2.0 {
			return errors.New("Presence penalty must be between -2.0 and 2.0")
		}
		oR.PresencePenalty = &presencePenalty

		return nil
	}
}

func WithSeed(seed int) completionOption {
	return func(oR *openaiRequest) error {
		oR.Seed = &seed

		return nil
	}
}

func WithStop(stop []string) completionOption {
	return func(oR *openaiRequest) error {
		oR.Stop = stop

		return nil
	}
}

func WithTemperature(temperature float64) completionOption {
	return func(oR *openaiRequest) error {
		if temperature < 0.0 || temperature > 2.0 {
			return errors.New("Temperature must be between 0.0 and 2.0.")
		}
		oR.Temperature = &temperature

		return nil
	}
}
func WithTopP(topP float64) completionOption {
	return func(oR *openaiRequest) error {
		if topP < 0.0 || topP > 1 {
			return errors.New("Top P must be between 0 and 1.")
		}
		oR.TopP = &topP

		return nil
	}
}
func WithUser(user string) completionOption {
	return func(oR *openaiRequest) error {
		if user == "" {
			return errors.New("Cannot set user to empty string.")
		}
		oR.User = user

		return nil
	}
}