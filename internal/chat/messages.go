package chat

import (
	"google.golang.org/genai"
)

// all this is kind of redundant rn, but it will make sense eventually when we will have multiple providers in place.
type Part *genai.Part

type Event string

const (
	ToolCall          Event = "TOOL CALL"
	ToolCallCompleted Event = "TOOL CALL COMPLETED"
	ResponseCompleted Event = "RESPONSE COMPLETED"
	ToolCallError     Event = "TOOL CALL ERROR"
	UserMessage       Event = "USER MESSAGE"
)

// TODO: make this as generic as possible to support all the possible events. if you want usage metadata, this is where you can add all that data.
// NOTE: this is also consumed by the history bubble to render the history view on the terminal.
type Message struct {
	Event Event
}
