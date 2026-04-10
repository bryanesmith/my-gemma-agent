package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"my-gemma-agent/agent"
)

func main() {
	fmt.Println("Gemma Agent Loop - Type 'quit' to exit")
	fmt.Println("Available tools: get_current_date")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	var messages []agent.Message

	systemMsg := agent.Message{
		Role: "system",
		Content: `You are a helpful assistant with access to tools. 
If you need to use a tool, respond with JSON: {"tool_calls": [{"function": "tool_name", "args": {}}]}
If tool results are provided in the conversation, respond naturally with the answer.`,
	}
	messages = append(messages, systemMsg)

	for {
		fmt.Print("You: ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "quit" || input == "exit" {
			break
		}

		if input == "" {
			continue
		}

		messages = append(messages, agent.Message{Role: "user", Content: input})

		fmt.Print("Gemma: ")
		response, err := agent.Chat(messages, true)
		if err != nil {
			log.Printf("Error calling Gemma: %v", err)
			continue
		}
		fmt.Println()
		messages = append(messages, *response)

		for len(response.ToolCalls) > 0 {
			for _, tc := range response.ToolCalls {
				result, err := agent.ExecuteTool(tc.Function.Name, tc.Function.Arguments)
				if err != nil {
					log.Printf("Error executing tool %s: %v", tc.Function.Name, err)
					continue
				}
				toolMsg := agent.Message{
					Role:       "tool",
					Content:    result,
					ToolCallID: tc.ID,
				}
				messages = append(messages, toolMsg)
			}

			messages = append(messages, agent.Message{
				Role:    "user",
				Content: "Based on the tool results above, please provide a natural language answer.",
			})

			fmt.Print("Gemma: ")
			response, err = agent.Chat(messages, false)
			if err != nil {
				log.Printf("Error calling Gemma: %v", err)
				break
			}
			fmt.Println()
			messages = append(messages, *response)
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Input error: %v", err)
			break
		}
	}

	fmt.Println("\nGoodbye!")
}
