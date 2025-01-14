package tool

import (
	"encoding/json"
)

type TaxiTool struct {
}

func (tx *TaxiTool) Definition() *ToolDefinition {
	return &ToolDefinition{
		Name:        "book_taxi",
		Description: "Book a taxi to go from source to destination",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"source": map[string]string{
					"description": "The starting location of the taxi ride",
					"type":        "string",
				},
				"destination": map[string]string{
					"description": "The final destination of the taxi ride",
					"type":        "string",
				},
			},
			"required": []string{
				"source",
				"destination",
			},
		},
	}
}
func (tx *TaxiTool) Run(arguments string) *WorkDone {
	var args struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
	}
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return &WorkDone{
			Error:    err,
			Response: "",
		}
	}

	ride := "booking taxi from " + args.Source + " to " + args.Destination

	return &WorkDone{
		Error:    nil,
		Response: ride,
	}
}
