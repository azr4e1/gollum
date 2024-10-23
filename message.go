package gollum

import "errors"

const (
	System    = "system"
	Assistant = "assistant"
	User      = "user"
)

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

func (c *chat) History() []Message {
	if c.messages == nil {
		return []Message{}
	}
	return c.messages
}

func (c *chat) Clear() {
	clear(c.messages)
}

func (c *chat) Pop() (Message, error) {
	if c.messages == nil || len(c.messages) == 0 {
		return Message{}, errors.New("chat is empty.")
	}

	lastEl := c.messages[len(c.messages)-1]
	c.messages = c.messages[:len(c.messages)-1]

	return lastEl, nil
}

func (c *chat) Len() int {
	if c.messages == nil {
		return 0
	}
	return len(c.messages)
}

func (c *chat) SystemMessages() []Message {
	messages := []Message{}
	if c.messages == nil {
		return messages
	}
	for _, m := range c.messages {
		if m.Role == System {
			messages = append(messages, m)
		}
	}

	return messages
}

func (c *chat) UserMessages() []Message {
	messages := []Message{}
	if c.messages == nil {
		return messages
	}
	for _, m := range c.messages {
		if m.Role == User {
			messages = append(messages, m)
		}
	}

	return messages
}

func (c *chat) AssistantMessages() []Message {
	messages := []Message{}
	if c.messages == nil {
		return messages
	}
	for _, m := range c.messages {
		if m.Role == Assistant {
			messages = append(messages, m)
		}
	}

	return messages
}

func SystemMessage(content string) Message {
	return Message{Role: System, Content: content}
}

func UserMessage(content string) Message {
	return Message{Role: User, Content: content}
}

func AssistantMessage(content string) Message {
	return Message{Role: Assistant, Content: content}
}
