# protocol

Foundation types for the [TAU](https://github.com/tailored-agentic-units) agent ecosystem.

```
go get github.com/tailored-agentic-units/protocol
```

## Packages

| Package | Purpose |
|---------|---------|
| `protocol` | `Protocol` constants (Chat, Vision, Tools, Embeddings, Audio), `Message`, `Role`, `ToolCall` types |
| `config` | `AgentConfig`, `ClientConfig`, `ProviderConfig`, `ModelConfig` with JSON serialization and merge support |
| `response` | Unified `Response` with `ContentBlock` interface (`TextBlock`, `ToolUseBlock`), `StreamingResponse`, `EmbeddingsResponse` |
| `model` | Runtime `Model` type bridging configuration to domain |
| `streaming` | `StreamReader` interface and `StreamLine` type for streamed LLM responses |

## Usage

```go
import (
    "github.com/tailored-agentic-units/protocol"
    "github.com/tailored-agentic-units/protocol/config"
)

// Create messages
messages := []protocol.Message{
    protocol.SystemMessage("You are a helpful assistant."),
    protocol.UserMessage("Hello!"),
}

// Load agent configuration from JSON
cfg, err := config.LoadAgentConfig("agent.json")
```

## Dependencies

Standard library only.
