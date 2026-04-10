package agent

import "encoding/json"

type Tool struct {
	Type     string `json:"type"`
	Function struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Parameters  struct {
			Type       string                 `json:"type"`
			Properties map[string]interface{} `json:"properties"`
			Required   []string               `json:"required"`
		} `json:"parameters"`
	} `json:"function"`
}

type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	} `json:"function"`
}

type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
	Tools    []Tool    `json:"tools,omitempty"`
	Format   string    `json:"format,omitempty"`
}

type ChatResponse struct {
	Message Message `json:"message"`
	Done    bool    `json:"done"`
}
