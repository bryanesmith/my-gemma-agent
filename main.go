package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type ChatResponse struct {
	Message Message `json:"message"`
	Done    bool    `json:"done"`
}

func getCurrentDate() string {
	return time.Now().Format("2006-01-02")
}

func callOllamaStreaming(messages []Message) (string, error) {
	req := ChatRequest{
		Model:    "gemma4:e4b",
		Messages: messages,
		Stream:   true,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	resp, err := http.Post("http://localhost:11434/api/chat", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var fullResponse strings.Builder
	decoder := json.NewDecoder(resp.Body)

	for {
		var chunk ChatResponse
		if err := decoder.Decode(&chunk); err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		if chunk.Message.Content != "" {
			fmt.Print(chunk.Message.Content)
			fullResponse.WriteString(chunk.Message.Content)
		}

		if chunk.Done {
			break
		}
	}

	return fullResponse.String(), nil
}

func main() {
	fmt.Println("Gemma Agent Loop - Type 'quit' to exit")
	fmt.Println("Available tools: get_current_date")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	var messages []Message

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

		if strings.Contains(strings.ToLower(input), "date") || strings.Contains(strings.ToLower(input), "time") {
			fmt.Printf("Agent: Today's date is %s\n", getCurrentDate())
			continue
		}

		messages = append(messages, Message{Role: "user", Content: input})

		fmt.Print("Gemma: ")
		response, err := callOllamaStreaming(messages)
		if err != nil {
			log.Printf("Error calling Gemma: %v", err)
			continue
		}
		fmt.Println()
		messages = append(messages, Message{Role: "assistant", Content: response})

		if err := scanner.Err(); err != nil {
			log.Printf("Input error: %v", err)
			break
		}
	}

	fmt.Println("\nGoodbye!")
}
