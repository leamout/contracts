package ai

import (
	"testing"
	"time"
)

func TestAudioFrameDuration(t *testing.T) {
	frame := AudioFrame{
		Data: make([]byte, 640),
		Format: AudioFormat{
			Encoding:     AudioEncodingPCM16LE,
			SampleRateHz: 16000,
			Channels:     1,
		},
	}
	if got, want := frame.Duration(), 20*time.Millisecond; got != want {
		t.Fatalf("Duration() = %v, want %v", got, want)
	}
}

func TestAudioFrameValidateRejectsOddPCM16(t *testing.T) {
	frame := AudioFrame{
		Data: []byte{1, 2, 3},
		Format: AudioFormat{
			Encoding:     AudioEncodingPCM16LE,
			SampleRateHz: 16000,
			Channels:     1,
		},
	}
	if err := frame.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want error")
	}
}
