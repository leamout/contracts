package ai

import (
	"context"
	"time"
)

// RealtimeEventType identifies normalized events from an integrated realtime
// speech model.
type RealtimeEventType string

const (
	RealtimeEventSpeechStarted   RealtimeEventType = "speech.started"
	RealtimeEventSpeechStopped   RealtimeEventType = "speech.stopped"
	RealtimeEventTranscriptDelta RealtimeEventType = "transcript.delta"
	RealtimeEventTranscriptFinal RealtimeEventType = "transcript.final"
	RealtimeEventResponseStarted RealtimeEventType = "response.started"
	RealtimeEventResponseDelta   RealtimeEventType = "response.delta"
	RealtimeEventResponseStopped RealtimeEventType = "response.stopped"
	RealtimeEventToolCall        RealtimeEventType = "tool.call"
	RealtimeEventUsage           RealtimeEventType = "usage"
	RealtimeEventError           RealtimeEventType = "error"
)

// TranscriptEvent contains normalized transcript text.
type TranscriptEvent struct {
	Text string `json:"text"`
}

// ResponseEvent contains normalized assistant response text.
type ResponseEvent struct {
	Text string `json:"text"`
}

// Usage contains provider-reported token usage when available.
type Usage struct {
	InputTokens  int `json:"input_tokens,omitempty"`
	OutputTokens int `json:"output_tokens,omitempty"`
	TotalTokens  int `json:"total_tokens,omitempty"`
}

// Failure contains a normalized provider failure.
type Failure struct {
	Source   string `json:"source,omitempty"`
	Code     string `json:"code,omitempty"`
	Message  string `json:"message"`
	Terminal bool   `json:"terminal,omitempty"`
}

// RealtimeEvent is emitted by an integrated realtime provider session.
type RealtimeEvent struct {
	Type       RealtimeEventType `json:"type"`
	Transcript *TranscriptEvent  `json:"transcript,omitempty"`
	Response   *ResponseEvent    `json:"response,omitempty"`
	ToolCall   *ToolCall         `json:"tool_call,omitempty"`
	Usage      *Usage            `json:"usage,omitempty"`
	Failure    *Failure          `json:"failure,omitempty"`
	ProviderID string            `json:"provider_id,omitempty"`
	OccurredAt time.Time         `json:"occurred_at,omitempty"`
}

// RealtimeRequest contains immutable configuration for one integrated
// speech-to-speech session.
type RealtimeRequest struct {
	Runtime      Runtime          `json:"runtime"`
	InputFormat  AudioFormat      `json:"input_format"`
	OutputFormat AudioFormat      `json:"output_format"`
	Language     string           `json:"language,omitempty"`
	Instructions string           `json:"instructions,omitempty"`
	Voice        string           `json:"voice,omitempty"`
	Tools        []ToolDefinition `json:"tools,omitempty"`
}

// RealtimeStream is a normalized integrated speech model session.
type RealtimeStream interface {
	SendAudio(context.Context, AudioFrame) error
	Interrupt(context.Context) error
	SubmitToolResult(context.Context, ToolResult) error
	Audio() <-chan AudioFrame
	Events() <-chan RealtimeEvent
	Close(context.Context) error
}

// Realtime is implemented by integrated speech model provider adapters.
type Realtime interface {
	Provider
	StartRealtime(context.Context, RealtimeRequest) (RealtimeStream, error)
}
