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

func CallOllamaStreaming(messages []Message) (string, error) {
	req := ChatRequest{
		Model:    ModelName,
		Messages: messages,
		Stream:   true,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	spinner := NewSpinner()
	spinner.Start("thinking...")

	resp, err := http.Post(OllamaURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		spinner.Stop()
		return "", err
	}
	defer resp.Body.Close()

	var fullResponse strings.Builder
	decoder := json.NewDecoder(resp.Body)
	firstToken := true

	for {
		var chunk ChatResponse
		if err := decoder.Decode(&chunk); err != nil {
			if firstToken {
				spinner.Stop()
			}
			if err == io.EOF {
				break
			}
			return "", err
		}

		if chunk.Message.Content != "" {
			if firstToken {
				spinner.Stop()
				firstToken = false
			}
			fmt.Print(chunk.Message.Content)
			fullResponse.WriteString(chunk.Message.Content)
		}

		if chunk.Done {
			break
		}
	}

	return fullResponse.String(), nil
}
