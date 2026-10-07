package agent

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type ToolType string

const (
	ToolTypeBuiltin ToolType = "builtin"
	ToolTypeWebhook ToolType = "webhook"

	BuiltinHangupCall   = "hangup_call"
	BuiltinTransferCall = "transfer_call"
	BuiltinHoldCall     = "hold_call"
	BuiltinResumeCall   = "resume_call"
	BuiltinSendDTMF     = "send_dtmf"
)

type ToolDefinition struct {
	Type        ToolType        `json:"type"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"`
	EndpointURL string          `json:"endpoint_url,omitempty"`
	TimeoutMS   *int32          `json:"timeout_ms,omitempty"`
	Enabled     *bool           `json:"enabled,omitempty"`
}

func (t ToolDefinition) Validate() error {
	switch t.Type {
	case ToolTypeBuiltin, ToolTypeWebhook:
	default:
		return fmt.Errorf("unsupported tool type %q", t.Type)
	}
	if strings.TrimSpace(t.Name) == "" {
		return fmt.Errorf("tool name is required")
	}
	if strings.TrimSpace(t.Description) == "" {
		return fmt.Errorf("tool description is required")
	}
	if len(t.Parameters) != 0 {
		var object map[string]json.RawMessage
		if err := json.Unmarshal(t.Parameters, &object); err != nil || object == nil {
			return fmt.Errorf("tool parameters must be a JSON object")
		}
	}
	if t.TimeoutMS != nil && (*t.TimeoutMS < 100 || *t.TimeoutMS > 30000) {
		return fmt.Errorf("tool timeout_ms must be between 100 and 30000")
	}

	switch t.Type {
	case ToolTypeBuiltin:
		if strings.TrimSpace(t.EndpointURL) != "" {
			return fmt.Errorf("builtin tools cannot define endpoint_url")
		}
		if !IsBuiltinTool(t.Name) {
			return fmt.Errorf("unsupported builtin tool %q", t.Name)
		}
	case ToolTypeWebhook:
		if err := validatePortableEndpoint(t.EndpointURL); err != nil {
			return err
		}
	}
	return nil
}

func ValidateTools(tools []ToolDefinition) error {
	seen := make(map[string]struct{}, len(tools))
	for _, tool := range tools {
		if err := tool.Validate(); err != nil {
			return err
		}
		name := strings.TrimSpace(tool.Name)
		if _, exists := seen[name]; exists {
			return fmt.Errorf("duplicate tool %q", name)
		}
		seen[name] = struct{}{}
	}
	return nil
}

func IsBuiltinTool(name string) bool {
	switch strings.TrimSpace(name) {
	case BuiltinHangupCall, BuiltinTransferCall, BuiltinHoldCall, BuiltinResumeCall, BuiltinSendDTMF:
		return true
	default:
		return false
	}
}

func validatePortableEndpoint(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("webhook tools require endpoint_url")
	}
	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") && len(value) > 3 {
		return nil
	}
	parsed, err := url.ParseRequestURI(value)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" {
		return fmt.Errorf("webhook endpoint_url must be HTTPS or an environment placeholder")
	}
	return nil
}
