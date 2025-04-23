package gollum

import (
	oai "github.com/azr4e1/gollum/openai"
)

func (t Tool) ToOpenai() oai.OpenaiTool {
	name := t.Function.Name
	description := t.Function.Description
	required := t.Function.Parameters.Required
	args := []oai.ToolArgument{}
	for name, param := range t.Function.Parameters.Properties {
		ta := oai.ToolArgument{
			Name:        name,
			Type:        oai.JsonTypes(param.Type),
			Description: param.Description,
			Enum:        param.Enum,
		}
		args = append(args, ta)
	}

	return oai.NewTool(name, description, args, required)
}
