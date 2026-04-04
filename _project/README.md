# protocol

Foundational types for the TAU agent ecosystem. Zero external dependencies — pure Go types that all upper layers build on.

## Vision

A stable, minimal type foundation that defines the language spoken by all TAU libraries. Protocol constants, message types, configuration structures, a unified response model, and streaming interfaces. Changes here ripple through the entire ecosystem, so the API surface is kept deliberately small and stable.

## Core Premise

LLM APIs are message-array-based. The protocol library reflects this reality directly — no convenience abstractions that hide the underlying contract. Upper layers (format, provider, agent) build on these types without wrapping or re-defining them.

## Phases

| Phase | Focus Area | Version Target |
|-------|-----------|----------------|
| Phase 1 - Foundation | All core types: protocol constants, Message, config, unified Response model, streaming interfaces | v0.1.0 |

## Architecture

```
protocol (root)    — Protocol constants, Message, Role, ToolCall, ToolFunction
  config/          — AgentConfig, ClientConfig, ModelConfig, ProviderConfig, Duration
  response/        — Response (unified), ContentBlock, TextBlock, ToolUseBlock,
                     StreamingResponse, EmbeddingsResponse, TokenUsage
  model/           — Model runtime type (config → protocol bridge)
  streaming/       — StreamReader interface, StreamLine (interfaces only)
```

### Dependency Hierarchy

```
Level 0: root (protocol constants, Message, ToolCall)
Level 0: config/ (configuration types — no internal deps)
Level 1: response/ (no internal deps — pure types)
Level 2: model/ (imports config/, root for protocol bridge)
Level 2: streaming/ (no internal deps — pure interfaces)
```

## Source References

### Root Package (Protocol Constants + Message Types)

**Port from**: `~/tau/kernel/core/protocol/`

| Source File | Destination | Action |
|-------------|-------------|--------|
| `~/tau/kernel/core/protocol/protocol.go` (76 lines) | `protocol.go` | Port with minor changes |
| `~/tau/kernel/core/protocol/message.go` (63 lines) | `message.go` | Port with additions |
| `~/tau/kernel/core/protocol/tool.go` (11 lines) | **Do not port** | Tool type moves to tau/format as ToolDefinition |

**go-agents reference**: `~/code/go-agents/pkg/protocol/protocol.go` (72 lines), `~/code/go-agents/pkg/protocol/message.go` (22 lines)

**Deviations from kernel source**:
- `protocol/tool.go` (`Tool` struct) is NOT ported here — it moves to tau/format as `ToolDefinition`. The kernel's `protocol.Tool` was `{Name, Description, Parameters}` which is identical to go-agents `format.ToolDefinition`.
- Add convenience constructors: `UserMessage(string)`, `SystemMessage(string)`, `Messages(Role, string) []Message` — new helpers not in kernel or go-agents.

**Deviations from go-agents source**:
- go-agents `Message` has bare `string` Role and no tool fields. tau/protocol keeps the kernel's richer type: typed `Role`, `ToolCalls []ToolCall`, `ToolCallID string`. This is critical for multi-turn tool-calling workflows.
- go-agents has no `Audio` protocol constant. tau/protocol adds it as reserved.
- go-agents `Message.Content` is `any`. Kernel's is also `any`. No change.

### config/ Package

**Port from**: `~/tau/kernel/core/config/`

| Source File | Destination | Action |
|-------------|-------------|--------|
| `~/tau/kernel/core/config/agent.go` (86 lines) | `config/agent.go` | Port with addition |
| `~/tau/kernel/core/config/client.go` (80 lines) | `config/client.go` | Port intact |
| `~/tau/kernel/core/config/duration.go` (52 lines) | `config/duration.go` | Port intact |
| `~/tau/kernel/core/config/model.go` (61 lines) | `config/model.go` | Port intact |
| `~/tau/kernel/core/config/provider.go` (41 lines) | `config/provider.go` | Port intact |
| `~/tau/kernel/core/config/doc.go` (30 lines) | `config/doc.go` | Update package doc |

**go-agents reference**: `~/code/go-agents/pkg/config/` — structurally identical to kernel (same author, same lineage).

**Deviations from kernel source**:
- `AgentConfig` gains a `Format string` field (e.g., `"openai"`, `"converse"`). go-agents has this field; kernel does not. This tells the agent which wire format to use.
- All import paths change from `github.com/tailored-agentic-units/kernel/core/config` to `github.com/tailored-agentic-units/protocol/config`.

**Deviations from go-agents source**:
- go-agents `AgentConfig` has `Format string` — adopt this.
- Otherwise structurally identical. Same `Merge()` patterns, same `Default*Config()` functions.

### response/ Package

**REWRITE** — this is the most significant change in the extraction.

**Kernel source** (being replaced): `~/tau/kernel/core/response/`

| Kernel File | Status |
|-------------|--------|
| `~/tau/kernel/core/response/chat.go` (59 lines) — `ChatResponse`, `ParseChat()` | **Replaced** by unified Response |
| `~/tau/kernel/core/response/tools.go` (38 lines) — `ToolsResponse`, `ParseTools()` | **Replaced** by unified Response |
| `~/tau/kernel/core/response/streaming.go` (70 lines) — `StreamingChunk` | **Replaced** by StreamingResponse |
| `~/tau/kernel/core/response/embeddings.go` (30 lines) — `EmbeddingsResponse` | **Ported** (separate type preserved) |
| `~/tau/kernel/core/response/audio.go` (47 lines) — `AudioResponse` | **Not ported** in v0.1.0 (Audio deferred) |
| `~/tau/kernel/core/response/usage.go` (10 lines) — `TokenUsage` | **Ported** with field name changes |
| `~/tau/kernel/core/response/parse.go` (46 lines) — `Parse()`, `ParseStreamChunk()` | **Not ported** — parsing moves to tau/format |

**Build from go-agents**: `~/code/go-agents/pkg/response/`

| go-agents File | Destination | Action |
|----------------|-------------|--------|
| `~/code/go-agents/pkg/response/content.go` (23 lines) | `response/content.go` | Port intact |
| `~/code/go-agents/pkg/response/response.go` (32 lines) | `response/response.go` | Port intact |
| `~/code/go-agents/pkg/response/streaming.go` (21 lines) | `response/streaming.go` | Port intact |
| `~/code/go-agents/pkg/response/usage.go` (10 lines) | `response/usage.go` | Port intact |
| `~/code/go-agents/pkg/response/embeddings.go` (9 lines) | `response/embeddings.go` | Port intact |

**What changes from kernel**:
- `ChatResponse` (OpenAI `choices[].message` shape) → **Replaced** by `Response` with `[]ContentBlock`
- `ToolsResponse` (separate type for tool calls) → **Merged** into `Response` — tool calls are `ToolUseBlock` content blocks
- `StreamingChunk` (OpenAI delta shape) → **Replaced** by `StreamingResponse` with accumulated content blocks
- `TokenUsage.PromptTokens`/`CompletionTokens` → `InputTokens`/`OutputTokens` (provider-neutral naming)
- `ParseChat()`, `ParseTools()`, `ParseStreamChunk()` functions → **Removed** from protocol, moved to tau/format's Format.Parse/ParseStreamChunk methods
- `AudioResponse` → **Deferred** to v0.2.0

**What changes from go-agents**:
- No changes. go-agents response types are adopted as-is.

### model/ Package

**Port from**: `~/tau/kernel/core/model/model.go` (41 lines)

| Source File | Destination | Action |
|-------------|-------------|--------|
| `~/tau/kernel/core/model/model.go` | `model/model.go` | Port with import path changes |

**go-agents reference**: `~/code/go-agents/pkg/model/model.go` (41 lines) — structurally identical.

**Deviations**: Import path changes only (`kernel/core/config` → `protocol/config`, `kernel/core/protocol` → root `protocol`).

### streaming/ Package

**Build from go-agents**: `~/code/go-agents/pkg/streaming/streaming.go` (28 lines) — **interfaces only**

| go-agents File | Destination | Action |
|----------------|-------------|--------|
| `~/code/go-agents/pkg/streaming/streaming.go` | `streaming/streaming.go` | Port interfaces only |

**What is NOT ported here** (goes to tau/provider instead):
- `~/code/go-agents/pkg/streaming/sse.go` (75 lines) — SSE reader implementation → tau/provider
- `~/code/go-agents/pkg/streaming/eventstream.go` (87 lines) — EventStream reader implementation → tau/provider

**Deviations from go-agents**:
- Only the `StreamReader` interface and `StreamLine` struct are ported. The SSE and EventStream implementations stay in tau/provider because they carry external dependencies (AWS EventStream SDK).
- This split does not exist in go-agents (all streaming types are in one package).

**Deviations from kernel**:
- Kernel has no streaming abstraction at all — streaming is handled inline in providers via `<-chan any`. This is entirely new infrastructure from go-agents.

## Dependencies

None. Stdlib only.

## Integration Points

- **tau/format** imports protocol constants, Message, response types, streaming interfaces
- **tau/provider** imports protocol constants, config types, streaming interfaces
- **tau/agent** imports everything transitively
- **tau/kernel** imports protocol/Message for session history, response types for runtime loop
