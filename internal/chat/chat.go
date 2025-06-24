package chat

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// NOTE: this is what the leftpane should look like after refactoring. this gots to be a bubble.
type Chat struct {
	ChatId   string
	Title    string
	viewport viewport.Model
	History  History
}

func (c Chat) Init() tea.Cmd {}
func (c Chat) View() string {}
func (c Chat) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case 
	}
}
