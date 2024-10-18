package openai

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type llmChat struct {
	messages []message
}

func NewChat(m ...message) llmChat {
	return llmChat{messages: m}
}

func (c *llmChat) Add(m ...message) {
	c.messages = append(c.messages, m...)
}

func (c *llmChat) GetHistory() []message {
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
