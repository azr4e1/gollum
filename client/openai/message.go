package openai

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openaiChat struct {
	messages []message
}

func NewChat(m ...message) openaiChat {
	return openaiChat{messages: m}
}

func (c *openaiChat) Add(m ...message) {
	c.messages = append(c.messages, m...)
}

func (c *openaiChat) GetHistory() []message {
	if c.messages == nil {
		return nil
	}
	return c.messages
}

func SystemMessage(content string) message {
	return message{Role: "system", Content: content}
}
func UserMessage(content string) message {
	return message{Role: "user", Content: content}
}
func AssistantMessage(content string) message {
	return message{Role: "assistant", Content: content}
}
