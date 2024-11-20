package message

const (
	system    = "system"
	assistant = "assistant"
	user      = "user"
)

type Message struct {
	Role      string
	Content   string
	ToolCalls []toolCall
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

func UserMessage(content string) Message {
	return Message{Role: user, Content: content}
}

func AssistantMessage(content string) Message {
	return Message{Role: assistant, Content: content}
}
