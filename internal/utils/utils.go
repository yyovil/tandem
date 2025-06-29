package utils

import (
	"encoding/json"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
)

type (
	Type           string
	InfoType       int
	Status         string
	ClearStatusMsg struct{}
	InfoMsg        struct {
		Type InfoType
		Msg  string
		TTL  time.Duration
	}
)

const (
	Requesting    Status = "requesting"
	Streaming     Status = "streaming"
	ToolCall      Status = "tool_call"
	ToolCompleted Status = "tool_completed"
	Idle          Status = "idle"
	Error         Status = "error"
)

// OpenAPI 3.0 Specified type.
const (
	TypeUnspecified Type = "TYPE_UNSPECIFIED"
	TypeString      Type = "STRING"
	TypeNumber      Type = "NUMBER"
	TypeInteger     Type = "INTEGER"
	TypeBoolean     Type = "BOOLEAN"
	TypeArray       Type = "ARRAY"
	TypeObject      Type = "OBJECT"
	TypeNULL        Type = "NULL"
)

const (
	InfoTypeInfo InfoType = iota
	InfoTypeWarn
	InfoTypeError
)

func Wordwrap(content string, width int) string {
	// ADHD: these breakpoints are silly.
	var breakpoints string = " ,-"
	return ansi.Wordwrap(content, width, breakpoints)
}

// Clamp ensures that the value is within the specified min and max range.
func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func UnmarshalJSONToMap(data string) (result map[string]any, err error) {
	err = json.Unmarshal([]byte(data), &result)
	return result, err
}

func CmdHandler(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

func ReportError(err error) tea.Cmd {
	return CmdHandler(InfoMsg{
		Type: InfoTypeError,
		Msg:  err.Error(),
	})
}

func ReportWarn(warn string) tea.Cmd {
	return CmdHandler(InfoMsg{
		Type: InfoTypeWarn,
		Msg:  warn,
	})
}

func ReportInfo(info string) tea.Cmd {
	return CmdHandler(InfoMsg{
		Type: InfoTypeInfo,
		Msg:  info,
	})
}