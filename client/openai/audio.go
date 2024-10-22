package openai

type openaiVoice string
type audioFormat string

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

type audioRequest struct {
	Model  string      `json:"model"`
	Input  string      `json:"input"`
	Voice  openaiVoice `json:"voice"`
	Format audioFormat `json:"response_format,omitempty"`
	Speed  *int        `json:"speed,omitempty"`
}
