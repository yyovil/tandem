package messages

// RunEvent represents the possible events sent by the run() functions.
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
