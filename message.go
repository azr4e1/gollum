package gollum

type Message struct {
	Role    string
	Content string
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
