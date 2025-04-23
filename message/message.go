package message

import "encoding/json"

const (
	system    = "system"
	assistant = "assistant"
	user      = "user"
)

type Message struct {
	Role      string
	Content   string
	ToolCalls []ToolCall
}

type ToolCall struct {
	Id       string       `json:"id"`
	Type     string       `json:"type"`
	Function ToolCallFunc `json:"function"`
}

type ToolCallFunc struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

func UserMessage(content string) Message {
	return Message{Role: user, Content: content}
}

func AssistantMessage(content string) Message {
	return Message{Role: assistant, Content: content}
}
