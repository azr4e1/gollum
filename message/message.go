package message

const (
	System    = "system"
	Assistant = "assistant"
	User      = "user"
)

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

func SystemMessage(content string) Message {
	return Message{Role: System, Content: content}
}

func UserMessage(content string) Message {
	return Message{Role: User, Content: content}
}

func AssistantMessage(content string) Message {
	return Message{Role: Assistant, Content: content}
}
