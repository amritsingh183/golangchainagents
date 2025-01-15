package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	toolbox "github.com/amritsingh183/golangchainagents/pkg/toolbox"
	"github.com/openai/openai-go"

	"github.com/tmc/langchaingo/llms"
	langchainOpenAI "github.com/tmc/langchaingo/llms/openai"
)

type TaxiTool struct {
}

func (tx *TaxiTool) Definition() *toolbox.ToolDefinition {
	return &toolbox.ToolDefinition{
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
func (tx *TaxiTool) Call(ctx context.Context, arguments string) (string, error) {
	var args struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
	}
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", err
	}

	ride := "booking taxi from " + args.Source + " to " + args.Destination

	return ride, nil
}

type WeatherTool struct {
}

func (w *WeatherTool) Definition() *toolbox.ToolDefinition {
	return &toolbox.ToolDefinition{
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
func (w *WeatherTool) Call(ctx context.Context, arguments string) *toolbox.WorkDone {
	// Extract the location from the function call arguments
	var args struct {
		Location string `json:"location"`
		Unit     string `json:"unit"`
	}
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", err
	}

	// Simulate getting weather data
	weatherData := fmt.Sprintf("Weather in %s is Sunny, 25° %s", args.Location, args.Unit)

	return weatherData, nil
}

var systemPrompt = `You are an expert in composing functions. You are given a question and a set of possible functions. 
Based on the question, you will need to make one or more function/tool calls to achieve the purpose. 
If none of the functions can be used, point it out and refuse to answer. 
If the given question lacks the parameters required by the function, also point it out.

You have access to the following tools:
<tools>{{ tools }}</tools>

The output MUST strictly adhere to the following format, and NO other text MUST be included.
The example format is as follows. Please make sure the parameter type is correct. If no function call is needed, please make the tool calls an empty list '[]'.
<tool_call>[
{"name": "func_name1", "arguments": {"argument1": "value1", "argument2": "value2"}},

(more tool calls as required)
]</tool_call>`

var userQuery = "I am planning to visit New York City but I would like to visit there only if the weather there is good. What is the whether there? I have a small problem, I can understand only in celsius and not in fahrenheit, so let me know the temperature in the format i understand. If the whether is good book me a taxi to go there from california"

var modelName = openai.ChatModelGPT4oMini

func main() {
	ctx := context.Background()
	os.Setenv("OPENAI_API_KEY", "OPENAI_API_KEY")
	callOptions := []llms.CallOption{
		// llms.WithTemperature(0.9),
		// llms.WithTopP(0.65),
	}
	toolBox := toolbox.ToolBox{
		&WeatherTool{},
		&TaxiTool{},
	}
	llmTools := make([]llms.Tool, len(toolBox))
	for i, tl := range toolBox {
		toolDef := tl.Definition()
		llmTools[i] = llms.Tool{
			Type: "function",
			Function: &llms.FunctionDefinition{
				Name:        toolDef.Name,
				Description: toolDef.Description,
				Parameters:  toolDef.Parameters,
			},
		}
	}

	callOptions = append(callOptions, llms.WithTools(llmTools))

	llm, err := langchainOpenAI.New(langchainOpenAI.WithModel(modelName))
	if err != nil {
		log.Fatal(err)
	}

	// type may be text, image etc
	messageHistory := []llms.MessageContent{
		{
			Role: llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{
				llms.TextContent{
					Text: systemPrompt,
				},
			},
		},
		{
			Role: llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{
				llms.TextContent{
					Text: userQuery,
				},
			},
		},
	}
	resp, err := llm.GenerateContent(ctx, messageHistory, callOptions...)
	if err != nil {
		log.Fatal(err)
	}
	for _, choice := range resp.Choices {
		// logger.Log("Assistant message ", completion.Choices[0].Message.JSON)
		toolCalls := choice.ToolCalls

		// Abort early if there are no tool calls
		if len(toolCalls) == 0 {
			log.Fatal("No function call")
		}
		for _, toolCall := range choice.ToolCalls {
			toolResponse, err := toolBox.UseTool(ctx, toolCall.FunctionCall.Name, toolCall.FunctionCall.Arguments)
			if err != nil {
				log.Fatal(err)
			}
			assistantResponse := llms.MessageContent{
				Role: llms.ChatMessageTypeAI,
				Parts: []llms.ContentPart{
					llms.ToolCall{
						ID:   toolCall.ID,
						Type: toolCall.Type,
						FunctionCall: &llms.FunctionCall{
							Name:      toolCall.FunctionCall.Name,
							Arguments: toolCall.FunctionCall.Arguments,
						},
					},
				},
			}
			log.Printf("\n\n ToolCallResponse: %v", toolResponse)
			tool_result := llms.MessageContent{
				Role: llms.ChatMessageTypeTool,
				Parts: []llms.ContentPart{
					llms.ToolCallResponse{
						ToolCallID: toolCall.ID,
						Name:       toolCall.FunctionCall.Name,
						Content:    toolResponse,
					},
				},
			}

			// Append tool_use to messageHistory
			messageHistory = append(messageHistory, assistantResponse)
			// Append tool_result to messageHistory
			messageHistory = append(messageHistory, tool_result)

		}

		messageHistory = append(messageHistory, llms.MessageContent{
			Role: llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{
				llms.TextContent{
					Text: "Please formulate your final response now",
				},
			},
		})

		// Send query to the model again, this time with a history containing its
		// request to invoke a tool and our response to the tool call.
		log.Println("Querying with tool response...")
		resp, err = llm.GenerateContent(ctx, messageHistory, callOptions...)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("\n\n Final Response: %v", resp.Choices[0].Content)
		// populate ai response
		assistantResponse := llms.TextParts(llms.ChatMessageTypeAI, resp.Choices[0].Content)
		messageHistory = append(messageHistory, assistantResponse)
		log.Println(messageHistory)
	}

}
