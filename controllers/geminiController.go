package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
)

func TestGemini(c *fiber.Ctx) error {
	// genaiKey := os.Getenv("GOOGLE_API_KEY")
	// if genaiKey == "" {
	// 	log.Fatal("please set GOOGLE_API_KEY")
	// }

	type RequestBody struct {
		Question string `json:"question"`
	}

	body := new(RequestBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	genaiKey := ""

	ctx := context.Background()

	llm, err := googleai.New(ctx, googleai.WithAPIKey(genaiKey))
	if err != nil {
		log.Fatal(err)
	}

	// Start by sending an initial question about the weather to the model, adding
	// "available tools" that include a getCurrentWeather function.
	// Thoroughout this sample, messageHistory collects the conversation history
	// with the model - this context is needed to ensure tool calling works
	// properly.

	fmt.Println(body.Question)

	messageHistory := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, body.Question),
	}
	resp, err := llm.GenerateContent(ctx, messageHistory, llms.WithTools(availableToolsGemini))
	if err != nil {
		log.Fatal(err)
	}

	// Translate the model's response into a MessageContent element that can be
	// added to messageHistory.
	respchoice := resp.Choices[0]
	assistantResponse := llms.TextParts(llms.ChatMessageTypeAI, respchoice.Content)
	for _, tc := range respchoice.ToolCalls {
		assistantResponse.Parts = append(assistantResponse.Parts, tc)
	}
	messageHistory = append(messageHistory, assistantResponse)

	// "Execute" tool calls by calling requested function
	for _, tc := range respchoice.ToolCalls {
		switch tc.FunctionCall.Name {
		case "getCurrentWeather":
			var args struct {
				Location string `json:"location"`
			}
			if err := json.Unmarshal([]byte(tc.FunctionCall.Arguments), &args); err != nil {
				log.Fatal(err)
			}
			if strings.Contains(args.Location, "Indonesia") {
				toolResponse := llms.MessageContent{
					Role: llms.ChatMessageTypeTool,
					Parts: []llms.ContentPart{
						llms.ToolCallResponse{
							Name:    tc.FunctionCall.Name,
							Content: "64 and sunny",
						},
					},
				}
				messageHistory = append(messageHistory, toolResponse)
			}
		case "getCurrentOmzet":
			var args struct {
				Location string `json:"location"`
			}
			if err := json.Unmarshal([]byte(tc.FunctionCall.Arguments), &args); err != nil {
				log.Fatal(err)
			}
			if strings.Contains(args.Location, "Malang") {
				toolResponse := llms.MessageContent{
					Role: llms.ChatMessageTypeTool,
					Parts: []llms.ContentPart{
						llms.ToolCallResponse{
							Name:    tc.FunctionCall.Name,
							Content: "1 million and keep rising",
						},
					},
				}
				messageHistory = append(messageHistory, toolResponse)
			} else {
				toolResponse := llms.MessageContent{
					Role: llms.ChatMessageTypeTool,
					Parts: []llms.ContentPart{
						llms.ToolCallResponse{
							Name:    tc.FunctionCall.Name,
							Content: "2 million and stagnant",
						},
					},
				}
				messageHistory = append(messageHistory, toolResponse)
			}
		default:
			log.Fatalf("got unexpected function call: %v", tc.FunctionCall.Name)
		}
	}

	resp, err = llm.GenerateContent(ctx, messageHistory, llms.WithTools(availableToolsGemini))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(messageHistory)
	fmt.Println(resp.Choices[0])

	fmt.Println("Response after tool call:")
	b, _ := json.MarshalIndent(resp.Choices[0], " ", "  ")
	// fmt.Println(string(b))

	return c.JSON(string(b))
}

// availableToolsGemini simulates the tools/functions we're making available for
// the model.
var availableToolsGemini = []llms.Tool{
	{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        "getCurrentWeather",
			Description: "Get the current weather in a given location",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"location": map[string]any{
						"type":        "string",
						"description": "The city and state, e.g. San Francisco, CA",
					},
				},
				"required": []string{"location"},
			},
		},
	},
	{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        "getCurrentOmzet",
			Description: "Get the current omzet in a given location",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"location": map[string]any{
						"type":        "string",
						"description": "The branch office name, e.g. Malang, Pasuruan, Jombang",
					},
				},
				"required": []string{"location"},
			},
		},
	},
}
