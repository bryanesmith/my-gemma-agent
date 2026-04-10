package agent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func GetAvailableTools() []Tool {
	return []Tool{
		{
			Type: "function",
			Function: struct {
				Name        string `json:"name"`
				Description string `json:"description"`
				Parameters  struct {
					Type       string                 `json:"type"`
					Properties map[string]interface{} `json:"properties"`
					Required   []string               `json:"required"`
				} `json:"parameters"`
			}{
				Name:        "get_current_date",
				Description: "Get the current date in YYYY-MM-DD format",
				Parameters: struct {
					Type       string                 `json:"type"`
					Properties map[string]interface{} `json:"properties"`
					Required   []string               `json:"required"`
				}{
					Type:       "object",
					Properties: map[string]interface{}{},
					Required:   []string{},
				},
			},
		},
	}
}

func ExecuteTool(name string, args json.RawMessage) (string, error) {
	switch name {
	case "get_current_date":
		return time.Now().Format("2006-01-02"), nil
	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}

func ParseToolCalls(content string) ([]ToolCall, error) {
	content = strings.TrimSpace(content)

	var wrapper struct {
		ToolCalls []struct {
			Function string          `json:"function"`
			Args     json.RawMessage `json:"args"`
		} `json:"tool_calls"`
	}
	if err := json.Unmarshal([]byte(content), &wrapper); err == nil && len(wrapper.ToolCalls) > 0 {
		toolCalls := make([]ToolCall, len(wrapper.ToolCalls))
		for i, tc := range wrapper.ToolCalls {
			toolCalls[i] = ToolCall{
				ID:   fmt.Sprintf("call_%d", i),
				Type: "function",
				Function: struct {
					Name      string          `json:"name"`
					Arguments json.RawMessage `json:"arguments"`
				}{
					Name:      tc.Function,
					Arguments: tc.Args,
				},
			}
		}
		return toolCalls, nil
	}

	var toolCalls []ToolCall
	if err := json.Unmarshal([]byte(content), &toolCalls); err != nil {
		return nil, err
	}
	return toolCalls, nil
}

func ExtractJSON(content string) string {
	start := strings.Index(content, "{")
	if start == -1 {
		return ""
	}

	depth := 0
	for i := start; i < len(content); i++ {
		switch content[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return content[start : i+1]
			}
		}
	}
	return ""
}
