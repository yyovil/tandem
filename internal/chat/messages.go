package chat

import (
	"google.golang.org/genai"
)

// all this is kind of redundant rn, but it will make sense eventually when we will have multiple providers in place.
type Part *genai.Part

// we need to make all these Messages a bubble.
type AddUserMsg struct{}
type AddResponseMsg struct{}
type ToolCallMsg struct{}
type ToolCallCompletedMsg struct{}
type ResponseCompletedMsg struct{}
