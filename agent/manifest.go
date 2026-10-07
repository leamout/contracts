package agent

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/leamout/contracts/ai"
)

const SchemaVersion = 1

type Engine string

const (
	EngineComposable Engine = "composable"
	EngineRealtime   Engine = "realtime"
)

type InterruptionPolicy string

const (
	InterruptionAllow    InterruptionPolicy = "allow"
	InterruptionDisabled InterruptionPolicy = "disabled"
)

type RecordingPolicy string

const (
	RecordingNone RecordingPolicy = "none"
	RecordingAll  RecordingPolicy = "all"
)

type ProviderBinding struct {
	Role     ai.Kind         `json:"role"`
	Provider string          `json:"provider"`
	Config   json.RawMessage `json:"config,omitempty"`
}

type Manifest struct {
	SchemaVersion      int                `json:"schema_version"`
	Name               string             `json:"name"`
	Engine             Engine             `json:"engine"`
	Language           string             `json:"language,omitempty"`
	Instructions       string             `json:"instructions"`
	Voice              string             `json:"voice,omitempty"`
	Providers          []ProviderBinding  `json:"providers"`
	Tools              []string           `json:"tools,omitempty"`
	InterruptionPolicy InterruptionPolicy `json:"interruption_policy,omitempty"`
	RecordingPolicy    RecordingPolicy    `json:"recording_policy,omitempty"`
}

func (m Manifest) Validate() error {
	if m.SchemaVersion != SchemaVersion {
		return fmt.Errorf("unsupported agent schema_version %d", m.SchemaVersion)
	}
	if strings.TrimSpace(m.Name) == "" {
		return fmt.Errorf("agent name is required")
	}
	if strings.TrimSpace(m.Instructions) == "" {
		return fmt.Errorf("agent instructions are required")
	}
	if err := validateEngine(m.Engine); err != nil {
		return err
	}
	if err := validatePolicies(m.InterruptionPolicy, m.RecordingPolicy); err != nil {
		return err
	}
	if err := validateBindings(m.Engine, m.Providers); err != nil {
		return err
	}

	seenTools := make(map[string]struct{}, len(m.Tools))
	for _, tool := range m.Tools {
		name := strings.TrimSpace(tool)
		if name == "" {
			return fmt.Errorf("agent tool name is required")
		}
		if _, exists := seenTools[name]; exists {
			return fmt.Errorf("duplicate agent tool %q", name)
		}
		seenTools[name] = struct{}{}
	}
	return nil
}

func validateEngine(engine Engine) error {
	switch engine {
	case EngineComposable, EngineRealtime:
		return nil
	default:
		return fmt.Errorf("unsupported agent engine %q", engine)
	}
}

func validatePolicies(interruption InterruptionPolicy, recording RecordingPolicy) error {
	switch interruption {
	case "", InterruptionAllow, InterruptionDisabled:
	default:
		return fmt.Errorf("unsupported interruption_policy %q", interruption)
	}
	switch recording {
	case "", RecordingNone, RecordingAll:
	default:
		return fmt.Errorf("unsupported recording_policy %q", recording)
	}
	return nil
}

func validateBindings(engine Engine, bindings []ProviderBinding) error {
	expected := map[ai.Kind]struct{}{}
	switch engine {
	case EngineComposable:
		expected[ai.KindSTT] = struct{}{}
		expected[ai.KindLLM] = struct{}{}
		expected[ai.KindTTS] = struct{}{}
	case EngineRealtime:
		expected[ai.KindRealtime] = struct{}{}
	default:
		return validateEngine(engine)
	}

	seen := make(map[ai.Kind]struct{}, len(bindings))
	for _, binding := range bindings {
		if _, ok := expected[binding.Role]; !ok {
			return fmt.Errorf("provider role %q is not valid for %s engine", binding.Role, engine)
		}
		if _, exists := seen[binding.Role]; exists {
			return fmt.Errorf("duplicate provider role %q", binding.Role)
		}
		if strings.TrimSpace(binding.Provider) == "" {
			return fmt.Errorf("provider is required for role %q", binding.Role)
		}
		if len(binding.Config) != 0 {
			var object map[string]json.RawMessage
			if err := json.Unmarshal(binding.Config, &object); err != nil || object == nil {
				return fmt.Errorf("provider config for role %q must be a JSON object", binding.Role)
			}
		}
		seen[binding.Role] = struct{}{}
	}

	for role := range expected {
		if _, ok := seen[role]; !ok {
			return fmt.Errorf("agent requires provider role %q", role)
		}
	}
	return nil
}
