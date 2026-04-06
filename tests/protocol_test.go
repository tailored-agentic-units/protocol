package protocol_test

import (
	"encoding/json"
	"testing"

	"github.com/tailored-agentic-units/protocol"
)

func TestProtocol_Constants(t *testing.T) {
	tests := []struct {
		name     string
		protocol protocol.Protocol
		expected string
	}{
		{"Chat", protocol.Chat, "chat"},
		{"Vision", protocol.Vision, "vision"},
		{"Tools", protocol.Tools, "tools"},
		{"Embeddings", protocol.Embeddings, "embeddings"},
		{"Audio", protocol.Audio, "audio"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.protocol) != tt.expected {
				t.Errorf("got %s, want %s", string(tt.protocol), tt.expected)
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		name     string
		protocol string
		expected bool
	}{
		{"chat valid", "chat", true},
		{"vision valid", "vision", true},
		{"tools valid", "tools", true},
		{"embeddings valid", "embeddings", true},
		{"audio valid", "audio", true},
		{"invalid", "invalid", false},
		{"empty string", "", false},
		{"uppercase", "CHAT", false},
		{"mixed case", "Chat", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := protocol.IsValid(tt.protocol)
			if result != tt.expected {
				t.Errorf("IsValid(%q) = %v, want %v", tt.protocol, result, tt.expected)
			}
		})
	}
}

func TestValidProtocols(t *testing.T) {
	result := protocol.ValidProtocols()

	expected := []protocol.Protocol{
		protocol.Chat,
		protocol.Vision,
		protocol.Tools,
		protocol.Embeddings,
		protocol.Audio,
	}

	if len(result) != len(expected) {
		t.Fatalf("got %d protocols, want %d", len(result), len(expected))
	}

	for i, p := range expected {
		if result[i] != p {
			t.Errorf("index %d: got %s, want %s", i, result[i], p)
		}
	}
}

func TestProtocolStrings(t *testing.T) {
	result := protocol.ProtocolStrings()
	expected := "chat, vision, tools, embeddings, audio"

	if result != expected {
		t.Errorf("got %q, want %q", result, expected)
	}
}

func TestProtocol_SupportsStreaming(t *testing.T) {
	tests := []struct {
		name     string
		protocol protocol.Protocol
		expected bool
	}{
		{"Chat supports streaming", protocol.Chat, true},
		{"Vision supports streaming", protocol.Vision, true},
		{"Tools supports streaming", protocol.Tools, true},
		{"Embeddings does not support streaming", protocol.Embeddings, false},
		{"Audio does not support streaming", protocol.Audio, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.protocol.SupportsStreaming(); got != tt.expected {
				t.Errorf("SupportsStreaming() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewMessage_StringContent(t *testing.T) {
	msg := protocol.NewMessage(protocol.RoleUser, "Hello, world!")

	if msg.Role != protocol.RoleUser {
		t.Errorf("got role %q, want %q", msg.Role, protocol.RoleUser)
	}

	content, ok := msg.Content.(string)
	if !ok {
		t.Errorf("content is not a string")
	} else if content != "Hello, world!" {
		t.Errorf("got content %q, want %q", content, "Hello, world!")
	}
}

func TestNewMessage_StructuredContent(t *testing.T) {
	content := map[string]any{
		"type": "text",
		"text": "Hello",
	}

	msg := protocol.NewMessage(protocol.RoleAssistant, content)

	if msg.Role != protocol.RoleAssistant {
		t.Errorf("got role %q, want %q", msg.Role, protocol.RoleAssistant)
	}

	if _, ok := msg.Content.(map[string]any); !ok {
		t.Errorf("content is not a map")
	}
}

func TestRole_Constants(t *testing.T) {
	tests := []struct {
		name     string
		role     protocol.Role
		expected string
	}{
		{"system", protocol.RoleSystem, "system"},
		{"user", protocol.RoleUser, "user"},
		{"assistant", protocol.RoleAssistant, "assistant"},
		{"tool", protocol.RoleTool, "tool"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.role) != tt.expected {
				t.Errorf("got %q, want %q", string(tt.role), tt.expected)
			}
		})
	}
}

func TestNewMessage_Roles(t *testing.T) {
	tests := []struct {
		name string
		role protocol.Role
	}{
		{"user", protocol.RoleUser},
		{"assistant", protocol.RoleAssistant},
		{"system", protocol.RoleSystem},
		{"tool", protocol.RoleTool},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := protocol.NewMessage(tt.role, "content")
			if msg.Role != tt.role {
				t.Errorf("got role %q, want %q", msg.Role, tt.role)
			}
		})
	}
}

func TestMessage_ToolCallFields(t *testing.T) {
	msg := protocol.Message{
		Role:       protocol.RoleAssistant,
		Content:    "I'll call a tool",
		ToolCallID: "",
		ToolCalls: []protocol.ToolCall{
			{
				ID:   "call_123",
				Type: "function",
				Function: protocol.ToolFunction{
					Name:      "get_weather",
					Arguments: `{"location":"Boston"}`,
				},
			},
		},
	}

	if len(msg.ToolCalls) != 1 {
		t.Fatalf("got %d tool calls, want 1", len(msg.ToolCalls))
	}

	tc := msg.ToolCalls[0]
	if tc.ID != "call_123" {
		t.Errorf("got ID %q, want %q", tc.ID, "call_123")
	}

	if tc.Type != "function" {
		t.Errorf("got Type %q, want %q", tc.Type, "function")
	}

	if tc.Function.Name != "get_weather" {
		t.Errorf("got Function.Name %q, want %q", tc.Function.Name, "get_weather")
	}

	if tc.Function.Arguments != `{"location":"Boston"}` {
		t.Errorf("got Function.Arguments %q, want %q", tc.Function.Arguments, `{"location":"Boston"}`)
	}
}

func TestToolCall_MarshalUnmarshal(t *testing.T) {
	original := protocol.ToolCall{
		ID:   "call_abc",
		Type: "function",
		Function: protocol.ToolFunction{
			Name:      "search",
			Arguments: `{"query":"golang"}`,
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded protocol.ToolCall
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.ID != original.ID {
		t.Errorf("got ID %q, want %q", decoded.ID, original.ID)
	}

	if decoded.Type != original.Type {
		t.Errorf("got Type %q, want %q", decoded.Type, original.Type)
	}

	if decoded.Function.Name != original.Function.Name {
		t.Errorf("got Function.Name %q, want %q", decoded.Function.Name, original.Function.Name)
	}

	if decoded.Function.Arguments != original.Function.Arguments {
		t.Errorf("got Function.Arguments %q, want %q", decoded.Function.Arguments, original.Function.Arguments)
	}
}

func TestNewToolCall(t *testing.T) {
	tc := protocol.NewToolCall("call_xyz", "calculator", `{"expression":"2+2"}`)

	if tc.ID != "call_xyz" {
		t.Errorf("got ID %q, want %q", tc.ID, "call_xyz")
	}

	if tc.Type != "function" {
		t.Errorf("got Type %q, want %q", tc.Type, "function")
	}

	if tc.Function.Name != "calculator" {
		t.Errorf("got Function.Name %q, want %q", tc.Function.Name, "calculator")
	}

	if tc.Function.Arguments != `{"expression":"2+2"}` {
		t.Errorf("got Function.Arguments %q, want %q", tc.Function.Arguments, `{"expression":"2+2"}`)
	}
}

func TestInitMessages(t *testing.T) {
	msgs := protocol.InitMessages(protocol.RoleUser, "Hello")

	if len(msgs) != 1 {
		t.Fatalf("got %d messages, want 1", len(msgs))
	}

	if msgs[0].Role != protocol.RoleUser {
		t.Errorf("got role %q, want %q", msgs[0].Role, protocol.RoleUser)
	}

	content, ok := msgs[0].Content.(string)
	if !ok {
		t.Fatal("content is not a string")
	}

	if content != "Hello" {
		t.Errorf("got content %q, want %q", content, "Hello")
	}
}

func TestUserMessage(t *testing.T) {
	msg := protocol.UserMessage("Hello from user")

	if msg.Role != protocol.RoleUser {
		t.Errorf("got role %q, want %q", msg.Role, protocol.RoleUser)
	}

	content, ok := msg.Content.(string)
	if !ok {
		t.Fatal("content is not a string")
	}

	if content != "Hello from user" {
		t.Errorf("got content %q, want %q", content, "Hello from user")
	}
}

func TestSystemMessage(t *testing.T) {
	msg := protocol.SystemMessage("You are helpful")

	if msg.Role != protocol.RoleSystem {
		t.Errorf("got role %q, want %q", msg.Role, protocol.RoleSystem)
	}

	content, ok := msg.Content.(string)
	if !ok {
		t.Fatal("content is not a string")
	}

	if content != "You are helpful" {
		t.Errorf("got content %q, want %q", content, "You are helpful")
	}
}

func TestAssistantMessage(t *testing.T) {
	msg := protocol.AssistantMessage("I can help")

	if msg.Role != protocol.RoleAssistant {
		t.Errorf("got role %q, want %q", msg.Role, protocol.RoleAssistant)
	}

	content, ok := msg.Content.(string)
	if !ok {
		t.Fatal("content is not a string")
	}

	if content != "I can help" {
		t.Errorf("got content %q, want %q", content, "I can help")
	}
}

func TestToolMessage(t *testing.T) {
	msg := protocol.ToolMessage("call_123", "result data")

	if msg.Role != protocol.RoleTool {
		t.Errorf("got role %q, want %q", msg.Role, protocol.RoleTool)
	}

	if msg.ToolCallID != "call_123" {
		t.Errorf("got ToolCallID %q, want %q", msg.ToolCallID, "call_123")
	}

	content, ok := msg.Content.(string)
	if !ok {
		t.Fatal("content is not a string")
	}

	if content != "result data" {
		t.Errorf("got content %q, want %q", content, "result data")
	}
}
