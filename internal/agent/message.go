package agent

import "github.com/yyovil/tandem/internal/tools"

type Message struct {
	Role         Role
	Type         Type
	Part         Part
	Files        []Blob // NOTE: maybe we can support attachments through URIs in future.
	FinishReason FinishReason
	TokenCount   int32
}

type Role string

const (
	// NOTE: we only define RoleTool for tool Response
	RoleTool      Role = "tool"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type Type string

const (
	ToolCallMsg          Type = "tool_call"
	ToolCallErrorMsg     Type = "tool_call_error"
	ToolResponseMsg      Type = "tool_response"
	ResponseMsg          Type = "response"
	ResponseCompletedMsg Type = "response_completed"
	UserMessageMsg       Type = "user_message"
)

// NOTE: exactly one field within a Part should be set.
type Part struct {
	Text       string
	ToolCalls  []tools.ToolCall
	ToolResult []tools.ToolResponse
}

type Blob struct {
	Name     string
	Data     []byte
	MimeType string
}

// NOTE: these are some of the most common reasons why a response can finish. including unexpected ones.
type FinishReason string

const (
	BecauseStop             FinishReason = "stop"
	BecauseMaxTokens        FinishReason = "max_tokens"
	BecausePermissionDenied FinishReason = "permission_denied"
	BecauseUnknown          FinishReason = "unknown"
)
