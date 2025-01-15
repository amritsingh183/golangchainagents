package toolbox

import "context"

type ToolDefinition struct {
	Type string
	// Name of the function
	Name string
	// Description is a description of the function.
	Description string
	// Parameters is a list of parameters for the function.
	Parameters map[string]interface{}
}
type GetsWorkDone interface {
	Call(context.Context, string) (string, error)
	Definition() *ToolDefinition
}

type ToolBox []GetsWorkDone

func (tb *ToolBox) UseTool(ctx context.Context, toolName string, toolArgs string) (string, error) {
	for _, tool := range *tb {
		if tool.Definition().Name == toolName {
			return tool.Call(ctx, toolArgs)
		}
	}
	return "", nil
}
