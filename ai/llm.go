package ai

import "context"

// Role identifies the semantic role of a conversation message.
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

// Message is a provider-neutral conversation message.
type Message struct {
	Role       Role       `json:"role"`
	Content    string     `json:"content,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

// LLMEvent is a normalized incremental model response.
type LLMEvent struct {
	ResponseID    string `json:"response_id,omitempty"`
	TextDelta     string `json:"text_delta,omitempty"`
	ToolCallID    string `json:"tool_call_id,omitempty"`
	ToolIndex     int    `json:"tool_index,omitempty"`
	ToolName      string `json:"tool_name,omitempty"`
	ToolArguments []byte `json:"tool_arguments,omitempty"`
	InputTokens   int    `json:"input_tokens,omitempty"`
	OutputTokens  int    `json:"output_tokens,omitempty"`
	TotalTokens   int    `json:"total_tokens,omitempty"`
	Done          bool   `json:"done,omitempty"`
	Err           error  `json:"-"`
}

// LLMRequest contains one model generation request.
type LLMRequest struct {
	Runtime      Runtime          `json:"runtime"`
	Messages     []Message        `json:"messages"`
	Tools        []ToolDefinition `json:"tools,omitempty"`
	Instructions string           `json:"instructions,omitempty"`
}

// LLMStream exposes incremental generation events.
type LLMStream interface {
	Events() <-chan LLMEvent
	Close() error
}

// LLM is implemented by language model provider adapters.
type LLM interface {
	Provider
	Generate(context.Context, LLMRequest) (LLMStream, error)
}
