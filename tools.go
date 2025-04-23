package gollum

import (
	"fmt"
	"reflect"
)

type jsonTypes string

const (
	JSONObject  jsonTypes = "object"
	JSONString  jsonTypes = "string"
	JSONNumber  jsonTypes = "number"
	JSONInteger jsonTypes = "integer"
	JSONArray   jsonTypes = "array"
	JSONBoolean jsonTypes = "boolean"
	JSONNull    jsonTypes = "null"
)

type Tool struct {
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
	Type        jsonTypes `json:"type"`
	Description string    `json:"description,omitempty"`
	Enum        []string  `json:"enum,omitempty"`
}

type ToolArgument struct {
	Name        string
	Type        jsonTypes
	Description string
	Enum        []string
}

func NewTool(name, description string, args []ToolArgument, required []string) Tool {
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

	tool := Tool{
		Type:     "function",
		Function: function,
	}

	return tool
}

func GenerateArguments[T any]() []ToolArgument {
	args := []ToolArgument{}

	t := reflect.TypeFor[T]()
	if t.Kind() != reflect.Struct {
		return nil
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := field.Tag.Get("json")
		if name == "" {
			name = field.Name
		}
		description := field.Tag.Get("jsonschema_description")

		var jsonType jsonTypes
		switch field.Type.Kind() {
		case reflect.Bool:
			jsonType = JSONBoolean
		case reflect.Int | reflect.Int8 | reflect.Int16 | reflect.Int32 | reflect.Int64 | reflect.Uint | reflect.Uint8 | reflect.Uint16 | reflect.Uint32 | reflect.Uint64 | reflect.Uintptr:
			jsonType = JSONInteger
		case reflect.Float64:
			jsonType = JSONNumber
		case reflect.Array | reflect.Slice:
			jsonType = JSONArray
		case reflect.String:
			jsonType = JSONString
		case reflect.Struct:
			jsonType = JSONObject
		default:
			jsonType = JSONNull
		}

		arg := ToolArgument{
			Name:        name,
			Description: description,
			Type:        jsonType,
		}
		args = append(args, arg)
	}

	return args
}
