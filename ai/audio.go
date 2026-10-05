package ai

import (
	"fmt"
	"time"
)

// AudioEncoding identifies the representation of audio frame data.
type AudioEncoding string

const (
	AudioEncodingPCM16LE AudioEncoding = "pcm_s16le"
	AudioEncodingMuLaw   AudioEncoding = "mulaw"
	AudioEncodingALaw    AudioEncoding = "alaw"
)

// AudioFormat describes audio exchanged with speech providers.
type AudioFormat struct {
	Encoding     AudioEncoding `json:"encoding"`
	SampleRateHz int           `json:"sample_rate_hz"`
	Channels     int           `json:"channels"`
}

// Validate checks format invariants shared by provider adapters.
func (f AudioFormat) Validate() error {
	switch f.Encoding {
	case AudioEncodingPCM16LE, AudioEncodingMuLaw, AudioEncodingALaw:
	case "":
		return fmt.Errorf("audio encoding is required")
	default:
		return fmt.Errorf("unsupported audio encoding %q", f.Encoding)
	}
	if f.SampleRateHz <= 0 {
		return fmt.Errorf("audio sample rate must be positive")
	}
	if f.Channels <= 0 {
		return fmt.Errorf("audio channel count must be positive")
	}
	return nil
}

// AudioFrame contains one immutable audio frame.
type AudioFrame struct {
	Data       []byte      `json:"-"`
	Format     AudioFormat `json:"format"`
	CapturedAt time.Time   `json:"captured_at,omitempty"`
}

// Validate checks frame and format invariants.
func (f AudioFrame) Validate() error {
	if err := f.Format.Validate(); err != nil {
		return err
	}
	if len(f.Data) == 0 {
		return fmt.Errorf("audio frame is empty")
	}
	if f.Format.Encoding == AudioEncodingPCM16LE && len(f.Data)%2 != 0 {
		return fmt.Errorf("PCM16 audio frame has odd byte length %d", len(f.Data))
	}
	return nil
}

// Duration returns the duration represented by an uncompressed PCM16 or G.711 frame.
func (f AudioFrame) Duration() time.Duration {
	if f.Format.SampleRateHz <= 0 || f.Format.Channels <= 0 || len(f.Data) == 0 {
		return 0
	}

	bytesPerSample := 0
	switch f.Format.Encoding {
	case AudioEncodingPCM16LE:
		bytesPerSample = 2
	case AudioEncodingMuLaw, AudioEncodingALaw:
		bytesPerSample = 1
	default:
		return 0
	}

	samples := len(f.Data) / (bytesPerSample * f.Format.Channels)
	return time.Duration(samples) * time.Second / time.Duration(f.Format.SampleRateHz)
}
