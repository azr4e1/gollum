package openai

type JsonTypes string

const (
	JSONObject  JsonTypes = "object"
	JSONString  JsonTypes = "string"
	JSONNumber  JsonTypes = "number"
	JSONInteger JsonTypes = "integer"
	JSONArray   JsonTypes = "array"
	JSONBoolean JsonTypes = "boolean"
	JSONNull    JsonTypes = "null"
)

type OpenaiTool struct {
	Type     string       `json:"type"`
	Function functionTool `json:"function"`
}

type functionTool struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Parameters  *functionParameter `json:"parameters,omitempty"`
}

type functionParameter struct {
	Type       string                      `json:"type"`
	Properties map[string]functionArgument `json:"properties"`
	Required   []string                    `json:"required,omitempty"`
}

type functionArgument struct {
	Type        JsonTypes `json:"type"`
	Description string    `json:"description,omitempty"`
	Enum        []string  `json:"enum,omitempty"`
}

type ToolArgument struct {
	Name        string
	Type        JsonTypes
	Description string
	Enum        []string
}

func NewTool(name, description string, args []ToolArgument, required []string) OpenaiTool {
	properties := make(map[string]functionArgument)
	for _, arg := range args {
		properties[arg.Name] = functionArgument{
			Type:        arg.Type,
			Description: arg.Description,
			Enum:        arg.Enum,
		}
	}
	parameters := &functionParameter{
		Type:       string(JSONObject),
		Properties: properties,
		Required:   required,
	}
	function := functionTool{
		Name:        name,
		Description: description,
		Parameters:  parameters,
	}

	tool := OpenaiTool{
		Type:     "function",
		Function: function,
	}

	return tool
}
