package openai

type openaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openaiChat struct {
	messages []openaiMessage
}

func NewChat(m ...openaiMessage) openaiChat {
	return openaiChat{messages: m}
}

func (c *openaiChat) Add(m ...openaiMessage) {
	c.messages = append(c.messages, m...)
}

func (c *openaiChat) GetHistory() []openaiMessage {
	if c.messages == nil {
		return []openaiMessage{}
	}
	return c.messages
}

func SystemMessage(content string) openaiMessage {
	return openaiMessage{Role: "system", Content: content}
}
func UserMessage(content string) openaiMessage {
	return openaiMessage{Role: "user", Content: content}
}
func AssistantMessage(content string) openaiMessage {
	return openaiMessage{Role: "assistant", Content: content}
}
