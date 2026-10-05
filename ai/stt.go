package ai

import "context"

// STTEventType identifies normalized speech-to-text stream events.
type STTEventType string

const (
	STTEventSpeechStarted   STTEventType = "speech.started"
	STTEventSpeechStopped   STTEventType = "speech.stopped"
	STTEventTranscriptDelta STTEventType = "transcript.delta"
	STTEventTranscriptFinal STTEventType = "transcript.final"
	STTEventError           STTEventType = "error"
)

// STTEvent is emitted by a streaming transcriber.
type STTEvent struct {
	Type       STTEventType `json:"type"`
	Text       string       `json:"text,omitempty"`
	ProviderID string       `json:"provider_id,omitempty"`
	Err        error        `json:"-"`
}

// STTRequest contains immutable configuration for one transcription stream.
type STTRequest struct {
	Runtime  Runtime     `json:"runtime"`
	Format   AudioFormat `json:"format"`
	Language string      `json:"language,omitempty"`
}

// STTStream is a bidirectional streaming speech-to-text session.
type STTStream interface {
	SendAudio(context.Context, AudioFrame) error
	Finalize(context.Context) error
	Events() <-chan STTEvent
	Close(context.Context) error
}

// STT is implemented by speech-to-text provider adapters.
type STT interface {
	Provider
	StartSTT(context.Context, STTRequest) (STTStream, error)
}
