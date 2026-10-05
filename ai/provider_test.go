package ai

import "testing"

func TestDescriptorValidate(t *testing.T) {
	tests := []struct {
		name    string
		value   Descriptor
		wantErr bool
	}{
		{
			name: "valid",
			value: Descriptor{
				ID:   "deepgram",
				Name: "Deepgram",
				Kind: KindSTT,
				Capabilities: []Capability{
					CapabilityStreaming,
					CapabilityTurnDetection,
				},
			},
		},
		{name: "missing id", value: Descriptor{Name: "Deepgram", Kind: KindSTT}, wantErr: true},
		{name: "missing name", value: Descriptor{ID: "deepgram", Kind: KindSTT}, wantErr: true},
		{name: "unsupported kind", value: Descriptor{ID: "deepgram", Name: "Deepgram", Kind: "other"}, wantErr: true},
		{
			name: "duplicate capability",
			value: Descriptor{
				ID:           "deepgram",
				Name:         "Deepgram",
				Kind:         KindSTT,
				Capabilities: []Capability{CapabilityStreaming, CapabilityStreaming},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.value.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDescriptorSupports(t *testing.T) {
	d := Descriptor{Capabilities: []Capability{CapabilityStreaming, CapabilityUsage}}
	if !d.Supports(CapabilityStreaming) {
		t.Fatal("Supports() = false, want true")
	}
	if d.Supports(CapabilityBargeIn) {
		t.Fatal("Supports() = true, want false")
	}
}
