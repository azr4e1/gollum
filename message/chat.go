package message

import (
	"errors"
)

type Chat struct {
	systemMessage Message
	messages      []Message
	limit         int
}

func (c *Chat) SetLimit(limit int) error {
	if limit < 0 {
		return errors.New("limit cannot be < 0.")
	}
	c.limit = limit
	return nil
}

func (c *Chat) SetSystemMessage(message string) {
	c.systemMessage = Message{role: "system", content: message}
}

func (c *Chat) UnsetSystemMessage() {
	c.systemMessage = Message{}
}

func NewChat(m ...Message) Chat {
	return Chat{messages: m}
}

func (c *Chat) Add(m ...Message) {
	if c.messages == nil {
		c.messages = []Message{}
	}

	c.messages = append(c.messages, m...)
	if c.limit > 0 && len(c.messages) > c.limit {
		trimmedMess := c.messages[len(c.messages)-c.limit:]
		c.messages = trimmedMess
	}
}

func (c Chat) History() []Message {
	if c.messages == nil {
		return []Message{}
	}
	return c.messages
}

func (c Chat) IsEmpty() bool {
	if c.messages == nil || len(c.messages) == 0 {
		return true
	}

	return false
}

func (c *Chat) Clear() {
	clear(c.messages)
}

func (c *Chat) Pop() (Message, error) {
	if c.messages == nil || len(c.messages) == 0 {
		return Message{}, errors.New("chat is empty.")
	}

	lastEl := c.messages[len(c.messages)-1]
	c.messages = c.messages[:len(c.messages)-1]

	return lastEl, nil
}

func (c *Chat) Len() int {
	if c.messages == nil {
		return 0
	}
	return len(c.messages)
}

func (c *Chat) SystemMessage() Message {
	return c.systemMessage
}

func (c *Chat) UserMessages() []Message {
	messages := []Message{}
	if c.messages == nil {
		return messages
	}
	for _, m := range c.messages {
		if m.Role() == user {
			messages = append(messages, m)
		}
	}

	return messages
}

func (c *Chat) AssistantMessages() []Message {
	messages := []Message{}
	if c.messages == nil {
		return messages
	}
	for _, m := range c.messages {
		if m.Role() == assistant {
			messages = append(messages, m)
		}
	}

	return messages
}
