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

// ToolRegistry manages available tools
type ToolRegistry struct {
	tools map[string]func() BaseTool
}

var globalRegistry = &ToolRegistry{
	tools: make(map[string]func() BaseTool),
}

// RegisterTool registers a tool constructor with the global registry
func RegisterTool(name string, constructor func() BaseTool) {
	globalRegistry.tools[name] = constructor
}

// GetTool creates a tool instance by name
func GetTool(name string) BaseTool {
	if constructor, exists := globalRegistry.tools[name]; exists {
		return constructor()
	}
	return nil
}

// GetToolsForAgent returns the tools configured for a specific agent
func GetToolsForAgent(toolNames []string) []BaseTool {
	var tools []BaseTool
	for _, name := range toolNames {
		if tool := GetTool(name); tool != nil {
			tools = append(tools, tool)
		}
	}
	return tools
}

// InitializeTools registers all available tools
func InitializeTools() {
	RegisterTool(VhsToolName, NewVhsTool)
	RegisterTool(FreezeToolName, NewFreezeTool)
	RegisterTool(GitAnalysisToolName, NewGitAnalysisTool)
}
