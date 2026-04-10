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

		if agent.IsDateQuery(input) {
			fmt.Printf("Agent: Today's date is %s\n", agent.GetCurrentDate())
			continue
		}

		messages = append(messages, agent.Message{Role: "user", Content: input})

		fmt.Print("Gemma: ")
		response, err := agent.ChatStreaming(messages)
		if err != nil {
			log.Printf("Error calling Gemma: %v", err)
			continue
		}
		fmt.Println()
		messages = append(messages, agent.Message{Role: "assistant", Content: response})

		if err := scanner.Err(); err != nil {
			log.Printf("Input error: %v", err)
			break
		}
	}

	fmt.Println("\nGoodbye!")
}
