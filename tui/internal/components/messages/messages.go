package messages

type RunEvent string

const (
	RunStarted            RunEvent = "RunStarted"
	RunResponseContent    RunEvent = "RunResponseContent"
	RunCompleted          RunEvent = "RunCompleted"
	RunError              RunEvent = "RunError"
	RunCancelled          RunEvent = "RunCancelled"
	RunPaused             RunEvent = "RunPaused"
	RunContinued          RunEvent = "RunContinued"
	ToolCallStarted       RunEvent = "ToolCallStarted"
	ToolCallCompleted     RunEvent = "ToolCallCompleted"
	ReasoningStarted      RunEvent = "ReasoningStarted"
	ReasoningStep         RunEvent = "ReasoningStep"
	ReasoningCompleted    RunEvent = "ReasoningCompleted"
	MemoryUpdateStarted   RunEvent = "MemoryUpdateStarted"
	MemoryUpdateCompleted RunEvent = "MemoryUpdateCompleted"
)

type ToolExecution struct {
	ToolCallId        string `json:"tool_call_id"`
	ToolName          string `json:"tool_name"`
	ToolArgs          any    `json:"tool_args"`
	ToolCallError     any    `json:"tool_call_error"`
	Result            string `json:"result"`
	Metrics           any    `json:"metrics"`
	StopAfterToolCall bool   `json:"stop_after_tool_call"`
	CreatedAt         int64  `json:"created_at"`
}

type RunResponse struct {
	CreatedAt   int64         `json:"created_at"`
	Event       RunEvent      `json:"event"`
	AgentId     string        `json:"agent_id"`
	RunId       string        `json:"run_id"`
	SessionId   string        `json:"session_id"`
	Content     string        `json:"content"`
	ContentType string        `json:"content_type"`
	Thinking    string        `json:"thinking"`
	Tool        ToolExecution `json:"tool"`
}
