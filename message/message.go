package message

import "encoding/json"

const (
	system    = "system"
	assistant = "assistant"
	user      = "user"
)

type Message struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type ToolCall struct {
	Id        string          `json:"id"`
	Type      string          `json:"type"`
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

func UserMessage(content string) Message {
	return Message{Role: user, Content: content}
}

func AssistantMessage(content string) Message {
	return Message{Role: assistant, Content: content}
}
