package ai

import "context"

// TextChunk is a unit of streamed text sent to a speech synthesizer.
type TextChunk struct {
	Text  string `json:"text"`
	Final bool   `json:"final,omitempty"`
}

// TTSEvent is emitted by a streaming text-to-speech provider.
type TTSEvent struct {
	Audio      AudioFrame `json:"audio"`
	ProviderID string     `json:"provider_id,omitempty"`
	Done       bool       `json:"done,omitempty"`
	Err        error      `json:"-"`
}

// TTSRequest contains immutable configuration for one synthesis stream.
type TTSRequest struct {
	Runtime  Runtime     `json:"runtime"`
	Format   AudioFormat `json:"format"`
	Voice    string      `json:"voice,omitempty"`
	Language string      `json:"language,omitempty"`
}

// TTSStream is a streaming text-to-speech session.
type TTSStream interface {
	SendText(context.Context, TextChunk) error
	Events() <-chan TTSEvent
	Close() error
}

// TTS is implemented by text-to-speech provider adapters.
type TTS interface {
	Provider
	StartTTS(context.Context, TTSRequest) (TTSStream, error)
}
