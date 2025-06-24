package chat

import "github.com/yyovil/tandem/internal/agent/tools"

// NOTE: this is also consumed by the history bubble to render the history view on the terminal.
type Message struct {
	Type         Type
	Part         Part
	Attachment   Blob
	FinishReason FinishReason
	TokenCount   int32
}

type Type string

const (
	ToolCallMsg          Type = "tool_call"
	ToolCallCompletedMsg Type = "tool_call_completed"
	ResponseCompletedMsg Type = "response_completed"
	UserMessageMsg       Type = "user_message"
)

type Part struct {
	Text     string
	ToolCall ToolCall
}

type ToolCall struct {
	Name tools.ToolName
	Args map[string]any
}

type Blob struct {
	Name     string
	Data     []byte
	MimeType string
}

// NOTE: these are some of the most common reasons why a response can finish. including unexpected ones.
type FinishReason string

const (
	BecauseEndTurn          FinishReason = "end_turn"
	BecauseMaxTokens        FinishReason = "max_tokens"
	BecauseToolUse          FinishReason = "tool_use"
	BecauseCanceled         FinishReason = "canceled"
	BecauseError            FinishReason = "error"
	BecausePermissionDenied FinishReason = "permission_denied"
	BecauseUnknown          FinishReason = "unknown"
)
