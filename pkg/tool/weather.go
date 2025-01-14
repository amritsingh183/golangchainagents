package tool

import (
	"encoding/json"
	"fmt"
)

type WeatherTool struct {
}

func (w *WeatherTool) Definition() *ToolDefinition {
	return &ToolDefinition{
		Name:        "get_weather",
		Description: "Get the weather in a given location",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"location": map[string]string{
					"description": "The name/location of place",
					"type":        "string",
				},
				"unit": map[string]interface{}{
					"enum": []string{
						"celsius",
						"fahrenheit",
					},
					"type": "string",
				},
			},
			"required": []string{
				"location",
			},
		},
	}
}
func (w *WeatherTool) Run(arguments string) *WorkDone {
	// Extract the location from the function call arguments
	var args struct {
		Location string `json:"location"`
		Unit     string `json:"unit"`
	}
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return &WorkDone{
			Error:    err,
			Response: "",
		}
	}

	// Simulate getting weather data
	weatherData := fmt.Sprintf("Weather in %s is Sunny, 25° %s", args.Location, args.Unit)

	return &WorkDone{
		Error:    nil,
		Response: weatherData,
	}
}
