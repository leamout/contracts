package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// Kind identifies the execution role implemented by a provider adapter.
type Kind string

const (
	KindSTT      Kind = "stt"
	KindLLM      Kind = "llm"
	KindTTS      Kind = "tts"
	KindRealtime Kind = "realtime"
)

// Capability describes optional behavior exposed by a provider adapter.
type Capability string

const (
	CapabilityStreaming     Capability = "streaming"
	CapabilityTurnDetection Capability = "turn_detection"
	CapabilityToolCalling   Capability = "tool_calling"
	CapabilityUsage         Capability = "usage"
	CapabilityBargeIn       Capability = "barge_in"
)

// Descriptor is the catalog metadata for one provider implementation.
// A vendor that supports multiple roles exposes one descriptor per Kind.
type Descriptor struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Kind         Kind         `json:"kind"`
	Capabilities []Capability `json:"capabilities,omitempty"`
}

// Validate checks the portable provider descriptor contract.
func (d Descriptor) Validate() error {
	if strings.TrimSpace(d.ID) == "" {
		return fmt.Errorf("provider id is required")
	}
	if strings.TrimSpace(d.Name) == "" {
		return fmt.Errorf("provider name is required")
	}

	switch d.Kind {
	case KindSTT, KindLLM, KindTTS, KindRealtime:
	default:
		return fmt.Errorf("unsupported provider kind %q", d.Kind)
	}

	seen := make(map[Capability]struct{}, len(d.Capabilities))
	for _, capability := range d.Capabilities {
		if strings.TrimSpace(string(capability)) == "" {
			return fmt.Errorf("provider capability is required")
		}
		if _, exists := seen[capability]; exists {
			return fmt.Errorf("duplicate provider capability %q", capability)
		}
		seen[capability] = struct{}{}
	}

	return nil
}

// Supports reports whether the descriptor advertises a capability.
func (d Descriptor) Supports(capability Capability) bool {
	for _, candidate := range d.Capabilities {
		if candidate == capability {
			return true
		}
	}
	return false
}

// Runtime contains the tenant-scoped opaque credential and provider-specific
// immutable configuration resolved before a provider session starts.
//
// Credential is intentionally opaque to the SDK. An adapter may interpret it
// as an API key, token, service-account document, or another provider-specific
// secret representation. Runtime values must never be persisted by adapters.
type Runtime struct {
	Credential string          `json:"-"`
	Config     json.RawMessage `json:"config,omitempty"`
}

// Provider is the common contract implemented by every adapter.
type Provider interface {
	Descriptor() Descriptor
}

// ConfigValidator is an optional provider contract for validating an adapter's
// provider-specific configuration before a runtime session is started.
type ConfigValidator interface {
	ValidateConfig(json.RawMessage) error
}

// CredentialVerifier is an optional provider contract for checking whether an
// opaque credential can authenticate with the upstream provider. Verification
// must not persist or return the supplied credential.
type CredentialVerifier interface {
	VerifyCredential(context.Context, string) error
}
