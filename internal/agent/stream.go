package agent

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Stream interface {
	Next() tea.Cmd
}

type StreamCreated struct {
	Stream <-chan Message
}
