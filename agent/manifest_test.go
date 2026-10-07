package agent

import (
	"encoding/json"
	"testing"

	"github.com/leamout/contracts/ai"
)

func TestManifestValidate(t *testing.T) {
	t.Parallel()
	manifest := Manifest{
		SchemaVersion: SchemaVersion,
		Name:          "Receptionist",
		Engine:        EngineComposable,
		Instructions:  "Answer calls.",
		Providers: []ProviderBinding{
			{Role: ai.KindSTT, Provider: "deepgram", Config: json.RawMessage(`{"model":"flux-general-en"}`)},
			{Role: ai.KindLLM, Provider: "groq"},
			{Role: ai.KindTTS, Provider: "cartesia"},
		},
		Tools: []string{"transfer_call"},
	}
	if err := manifest.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestManifestValidateRequiresEngineTopology(t *testing.T) {
	t.Parallel()
	manifest := Manifest{
		SchemaVersion: SchemaVersion,
		Name:          "Realtime",
		Engine:        EngineRealtime,
		Instructions:  "Answer calls.",
		Providers: []ProviderBinding{
			{Role: ai.KindSTT, Provider: "deepgram"},
		},
	}
	if err := manifest.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want topology error")
	}
}

func TestToolDefinitionValidate(t *testing.T) {
	t.Parallel()
	timeout := int32(5000)
	tool := ToolDefinition{
		Type:        ToolTypeWebhook,
		Name:        "lookup_customer",
		Description: "Look up a customer.",
		Parameters:  json.RawMessage(`{"type":"object"}`),
		EndpointURL: "${CUSTOMER_WEBHOOK_URL}",
		TimeoutMS:   &timeout,
	}
	if err := tool.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestRoutingManifestValidate(t *testing.T) {
	t.Parallel()
	manifest := RoutingManifest{
		SchemaVersion: SchemaVersion,
		Name:          "Business hours",
		DefaultAgent:  "receptionist",
		Routes: []Route{
			{Schedule: "after_hours", Agent: "after-hours"},
		},
		Schedules: map[string]Schedule{
			"after_hours": {
				Timezone: "UTC",
				Hours: []ScheduleHour{
					{Days: []string{"mon"}, Start: "17:00", End: "23:59"},
				},
			},
		},
	}
	if err := manifest.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if err := manifest.ValidateAgentReferences(map[string]struct{}{
		"receptionist": {},
		"after-hours":  {},
	}); err != nil {
		t.Fatalf("ValidateAgentReferences() error = %v", err)
	}
}
