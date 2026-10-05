# Leamout SDK

Shared SDK and extension contracts for Leamout runtimes, AI providers, carriers, and integrations.

## AI provider contracts

The `ai` package defines provider-neutral contracts for composable and integrated realtime voice agents. Provider implementations live outside the runtime and depend on these contracts rather than Leamout internals.

Supported roles:

- STT — streaming speech-to-text
- LLM — streaming language-model generation and tool calls
- TTS — streaming text-to-speech
- Realtime — integrated speech-to-speech providers

A vendor may implement more than one role by exposing a separate provider implementation for each `ai.Kind`.

```go
package example

import "github.com/leamout/sdk/ai"

type Provider struct{}

func (Provider) Descriptor() ai.Descriptor {
	return ai.Descriptor{
		ID:   "example",
		Name: "Example",
		Kind: ai.KindSTT,
		Capabilities: []ai.Capability{
			ai.CapabilityStreaming,
		},
	}
}
```

`ai.Runtime.Credential` is deliberately opaque. Provider adapters may interpret it as an API key, bearer token, service-account document, or another provider-specific secret representation. Provider-specific settings stay in `ai.Runtime.Config`.

Adapters may additionally implement `ai.ConfigValidator` to validate provider-specific configuration and `ai.CredentialVerifier` to verify upstream authentication without moving vendor logic into the Leamout runtime.

The shared audio contract supports PCM16 little-endian plus G.711 mu-law and A-law so speech adapters can work with both AI-native and telephony-native audio paths.

The SDK intentionally contains no vendor clients, credentials storage, runtime orchestration, or provider registry. Those concerns belong to `leamout/ai-providers` and the Leamout runtime.

## License

Apache License 2.0.
