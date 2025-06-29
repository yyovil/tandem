package agent

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Stream interface {
	// NOTE: pops the next message out of the Message Channel.
	Next() tea.Cmd
}

// NOTE: every time you send a message to the agent, a new stream is created.
type StreamCreated struct {
	Stream <-chan Message
}
