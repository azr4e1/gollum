package openai

type Message struct {
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
	messages []Message
}

func NewChat(m ...Message) chat {
	return chat{messages: m}
}

func (c *chat) Add(m ...Message) {
	c.messages = append(c.messages, m...)
}

func (c *chat) GetHistory() []Message {
	if c.messages == nil {
		return []Message{}
	}
	return c.messages
}

func SystemMessage(content string) Message {
	return Message{Role: "system", Content: content}
}
func UserMessage(content string) Message {
	return Message{Role: "user", Content: content}
}
func AssistantMessage(content string) Message {
	return Message{Role: "assistant", Content: content}
}
