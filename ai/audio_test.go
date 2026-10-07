package ai

import (
	"testing"
	"time"
)

func TestAudioFrameDuration(t *testing.T) {
	tests := []struct {
		name  string
		frame AudioFrame
		want  time.Duration
	}{
		{
			name: "pcm16 16kHz mono",
			frame: AudioFrame{
				Data: make([]byte, 640),
				Format: AudioFormat{
					Encoding:     AudioEncodingPCM16LE,
					SampleRateHz: 16000,
					Channels:     1,
				},
			},
			want: 20 * time.Millisecond,
		},
		{
			name: "mulaw 8kHz mono",
			frame: AudioFrame{
				Data: make([]byte, 160),
				Format: AudioFormat{
					Encoding:     AudioEncodingMuLaw,
					SampleRateHz: 8000,
					Channels:     1,
				},
			},
			want: 20 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.frame.Duration(); got != tt.want {
				t.Fatalf("Duration() = %v, want %v", got, tt.want)
			}
		})
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
