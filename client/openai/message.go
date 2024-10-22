package openai

type message struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []toolCall `json:"tool_calls,omitempty"`
}

type toolCall struct {
	Id       string       `json:"id"`
	Type     string       `json:"type"`
	Function toolCallFunc `json:"function"`
}

type toolCallFunc struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type chat struct {
	messages []message
}

func NewChat(m ...message) chat {
	return chat{messages: m}
}

func (c *chat) Add(m ...message) {
	c.messages = append(c.messages, m...)
}

func (c *chat) GetHistory() []message {
	if c.messages == nil {
		return []message{}
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
