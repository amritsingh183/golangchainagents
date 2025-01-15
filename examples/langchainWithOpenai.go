package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	toolbox "github.com/amritsingh183/golangchainagents/pkg/toolbox"

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
func (tx *TaxiTool) Call(ctx context.Context, arguments string) *toolbox.WorkDone {
	var args struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
	}
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return &toolbox.WorkDone{
			Error:    err,
			Response: "",
		}
	}

	ride := "booking taxi from " + args.Source + " to " + args.Destination

	return &toolbox.WorkDone{
		Error:    nil,
		Response: ride,
	}
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
		return &toolbox.WorkDone{
			Error:    err,
			Response: "",
		}
	}

	// Simulate getting weather data
	weatherData := fmt.Sprintf("Weather in %s is Sunny, 25° %s", args.Location, args.Unit)

	return &toolbox.WorkDone{
		Error:    nil,
		Response: weatherData,
	}
}

var templates = map[string]string{
	"smolM2": `You are an expert in composing functions. You are given a question and a set of possible functions. 
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
]</tool_call>`,
	"llama3": `You are a function calling AI model with epertise in in composing functions. You are given a question and a set of possible functions. You are provided with function signatures within <tools></tools> XML tags. Based on the question, you will need to make one or more function/tool calls to achieve the purpose. If none of the function can be used, point it out. If the given question lacks the parameters required by the function,also point it out. You should only return the function call in tools call sections.If you decide to invoke any of the function(s), you MUST put it in the format of [func_name1(params_name1=params_value1, params_name2=params_value2...)]\n
You SHOULD NOT include any other text in the response.\n`,
	"groq": `You are a function calling AI model. You are provided with function signatures within <tools></tools> XML tags. You may call one or more functions to assist with the user query. Don't make assumptions about what values to plug into functions. For each function call return a json object with function name and arguments within <tool_call></tool_call> XML tags as follows:
<tool_call>
{"name": <function-name>,"arguments": <args-dict>}
</tool_call>`,
}

func main() {
	var baseURL, apiKey, systemPrompt string
	os.Setenv("OPENAI_API_KEY", "OPENAI_API_KEY")
	models := map[string][]string{
		"SmolLM2-1.7B-Instruct-f16": {
			"http://127.0.0.1:11434/v1/", //baseurl
			templates["smolM2"],
			"y", //is local
			"",  //apikey
		},
		"Llama-3-Groq-8B-Tool-Use.Q8_0": {
			"http://127.0.0.1:11434/v1/",
			templates["groq"],
			"y",
			"",
		},
		"Llama-3.2-3B-Instruct-f16": {
			"http://127.0.0.1:11434/v1/",
			templates["llama3"],
			"y",
			"",
		},

		// openai.ChatModelGPT4oMini: {
		// 	"",
		// 	templates["llama3"],
		// 	"n",
		// 	os.Getenv("OPENAI_API_KEY"),
		// },
	}
	for modelName, modelMeta := range models {
		if modelMeta[2] == "y" {
			baseURL = modelMeta[0]
		} else {
			baseURL = ""
		}
		apiKey = modelMeta[3]
		systemPrompt = modelMeta[1]
		log.Println("Testing using langchain and ollama " + modelName)
		ctx := context.Background()
		makeLLMCall(ctx, baseURL, apiKey, modelName, systemPrompt)
	}

}
func makeLLMCall(ctx context.Context, baseURL, apiKey, modelName, systemPrompt string) {
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

	userQueries := []string{
		"I am planning to visit New York City but I would like to visit there only if the weather there is good. What is the whether there? I have a small problem, I can understand only in degree fahrenheit and not in celsius, so let me know the temperature in the format i understand",
		"I am planning to visit New York City but I would like to visit there only if the weather there is good. What is the whether there? I have a small problem, I can understand only in celsius and not in fahrenheit, so let me know the temperature in the format i understand",
		"I am planning to visit New York City but I would like to visit there only if the weather there is good. What is the whether there? I have a small problem, I can understand only in degree fahrenheit and not in celsius, so let me know the temperature in the format i understand. If the whether is good book me a taxi to go there from california",
		"I am planning to visit New York City but I would like to visit there only if the weather there is good. What is the whether there? I have a small problem, I can understand only in celsius and not in fahrenheit, so let me know the temperature in the format i understand. If the whether is good book me a taxi to go there from california",
		"I am planning to visit New York City by taxi but I would like to visit there only if the weather there is good. What is the whether there? I have a small problem, I can understand only in celsius and not in fahrenheit, so let me know the temperature in the format i understand. If the whether is good , i will go there tomorrow",
	}

	llm, err := langchainOpenAI.New(langchainOpenAI.WithBaseURL(baseURL), langchainOpenAI.WithModel(modelName))
	if err != nil {
		log.Fatal(err)
	}

	for _, userQuery := range userQueries {

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
				toolResponse := toolBox.UseTool(ctx, toolCall.FunctionCall.Name, toolCall.FunctionCall.Arguments)
				if toolResponse.Error != nil {
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
				log.Printf("\n\n ToolCallResponse: %v", toolResponse.Response)
				tool_result := llms.MessageContent{
					Role: llms.ChatMessageTypeTool,
					Parts: []llms.ContentPart{
						llms.ToolCallResponse{
							ToolCallID: toolCall.ID,
							Name:       toolCall.FunctionCall.Name,
							Content:    toolResponse.Response,
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

}
