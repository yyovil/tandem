package chat

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Stream interface {
	Next() tea.Cmd
}
