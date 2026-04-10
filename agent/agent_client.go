package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const OllamaURL = "http://localhost:11434/api/chat"
const ModelName = "gemma4:e4b"

func Chat(messages []Message, jsonMode bool) (*Message, error) {
	tools := GetAvailableTools()
	req := ChatRequest{
		Model:    ModelName,
		Messages: messages,
		Stream:   true,
		Tools:    tools,
	}
	if jsonMode {
		req.Format = "json"
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	spinner := NewSpinner()
	spinner.Start("thinking...")
	var firstToken = true
	defer func() {
		if firstToken {
			spinner.Stop()
		}
	}()

	resp, err := http.Post(OllamaURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var message Message
	message.Role = "assistant"

	var fullResponse strings.Builder
	decoder := json.NewDecoder(resp.Body)
	var firstTokenBuffer strings.Builder
	var isToolCall bool

	for {
		var chunk ChatResponse
		if err := decoder.Decode(&chunk); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if chunk.Message.Content != "" {
			if firstToken {
				spinner.Stop()
				firstToken = false
				firstTokenBuffer.WriteString(chunk.Message.Content)
				if strings.HasPrefix(strings.TrimSpace(firstTokenBuffer.String()), "{") {
					isToolCall = true
				} else {
					fmt.Print(firstTokenBuffer.String())
				}
				fullResponse.WriteString(chunk.Message.Content)
			} else if !isToolCall {
				fmt.Print(chunk.Message.Content)
				fullResponse.WriteString(chunk.Message.Content)
			} else {
				fullResponse.WriteString(chunk.Message.Content)
			}
		}

		if chunk.Message.ToolCalls != nil {
			message.ToolCalls = chunk.Message.ToolCalls
		}

		if chunk.Done {
			break
		}
	}

	message.Content = fullResponse.String()

	if message.Content != "" {
		toolCalls, err := ParseToolCalls(message.Content)
		if err == nil && len(toolCalls) > 0 {
			message.ToolCalls = toolCalls
			message.Content = ""
		}
	}

	return &message, nil
}
