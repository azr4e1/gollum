package gemini

const (
	System    = "system"
	Assistant = "model"
	User      = "user"
)

type Parts [](map[string]string)

type Message struct {
	Role string `json:"role"`
	Part Parts  `json:"parts"`
	// ToolCalls []toolCall `json:"tool_calls,omitempty"`
}
