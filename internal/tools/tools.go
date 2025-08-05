package tools

import (
	"context"
	"encoding/json"
)

type ToolInfo struct {
	Name        string
	Description string
	Parameters  map[string]any
	Required    []string
}

// ToolCallTracker tracks tool calls per session to enforce limits
type ToolCallTracker struct {
	sessionToolCalls map[string]map[string]int // sessionID -> toolName -> count
}

// NewToolCallTracker creates a new tool call tracker
func NewToolCallTracker() *ToolCallTracker {
	return &ToolCallTracker{
		sessionToolCalls: make(map[string]map[string]int),
	}
}

// IncrementToolCall increments the call count for a tool in a session
func (t *ToolCallTracker) IncrementToolCall(sessionID, toolName string) {
	if t.sessionToolCalls[sessionID] == nil {
		t.sessionToolCalls[sessionID] = make(map[string]int)
	}
	t.sessionToolCalls[sessionID][toolName]++
}

// GetToolCallCount returns the number of calls for a tool in a session
func (t *ToolCallTracker) GetToolCallCount(sessionID, toolName string) int {
	if t.sessionToolCalls[sessionID] == nil {
		return 0
	}
	return t.sessionToolCalls[sessionID][toolName]
}

// ResetSession clears the tool call counts for a session
func (t *ToolCallTracker) ResetSession(sessionID string) {
	delete(t.sessionToolCalls, sessionID)
}

type toolResponseType string

type (
	sessionIDContextKey string
	messageIDContextKey string
)

const (
	ToolResponseTypeText  toolResponseType = "text"
	ToolResponseTypeImage toolResponseType = "image"

	SessionIDContextKey sessionIDContextKey = "session_id"
	MessageIDContextKey messageIDContextKey = "message_id"
)

type ToolResponse struct {
	Type     toolResponseType `json:"type"`
	Content  string           `json:"content"`
	Metadata string           `json:"metadata,omitempty"`
	IsError  bool             `json:"is_error"`
}

func NewTextResponse(content string) ToolResponse {
	return ToolResponse{
		Type:    ToolResponseTypeText,
		Content: content,
	}
}

func WithResponseMetadata(response ToolResponse, metadata any) ToolResponse {
	if metadata != nil {
		metadataBytes, err := json.Marshal(metadata)
		if err != nil {
			return response
		}
		response.Metadata = string(metadataBytes)
	}
	return response
}

func NewTextErrorResponse(content string) ToolResponse {
	return ToolResponse{
		Type:    ToolResponseTypeText,
		Content: content,
		IsError: true,
	}
}

type ToolCall struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Input string `json:"input"`
}

type BaseTool interface {
	Info() ToolInfo
	Run(ctx context.Context, params ToolCall) (ToolResponse, error)
}

func GetContextValues(ctx context.Context) (string, string) {
	sessionID := ctx.Value(SessionIDContextKey)
	messageID := ctx.Value(MessageIDContextKey)
	if sessionID == nil {
		return "", ""
	}
	if messageID == nil {
		return sessionID.(string), ""
	}
	return sessionID.(string), messageID.(string)
}
