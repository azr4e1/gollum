package message

const (
	system    = "system"
	assistant = "assistant"
	user      = "user"
)

type Message struct {
	role      string
	content   string
	toolCalls []toolCall
}

func (m Message) Role() string {
	return m.role
}

func (m Message) Content() string {
	return m.content
}

func (m Message) Tools() []toolCall {
	return m.toolCalls
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
	return Message{role: user, content: content}
}

func AssistantMessage(content string) Message {
	return Message{role: assistant, content: content}
}
