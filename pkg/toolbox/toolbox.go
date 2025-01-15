package tool

type ToolDefinition struct {
	Type string
	// Name of the function
	Name string
	// Description is a description of the function.
	Description string
	// Parameters is a list of parameters for the function.
	Parameters map[string]interface{}
}
type WorkDone struct {
	Error    error
	Response string
}
type GetsWorkDone interface {
	// jsonString is {"argName":"argValue"..}
	Run(jsonString string) *WorkDone
	// name, description, parameters
	Definition() *ToolDefinition
}

type ToolBox []GetsWorkDone

func (tb *ToolBox) UseTool(toolName string, toolArgs string) *WorkDone {
	for _, tool := range *tb {
		if tool.Definition().Name == toolName {
			return tool.Run(toolArgs)
		}
	}
	return nil
}
