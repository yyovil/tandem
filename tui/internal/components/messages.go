package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type UserMessage struct {
	prompt         string
	attachmentName string
	width, height  int
}

// UserMsgAddedMsg is a message indicating a user message should be added to the leftpane viewport.
type UserMsgAddedMsg struct {
	UserMsg UserMessage
}

func (m *UserMessage) Init() tea.Cmd {
	return nil
}

func (m *UserMessage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *UserMessage) View() string {
	usermsgStyle := lipgloss.
		NewStyle().
		Border(lipgloss.InnerHalfBlockBorder(), false, false, false, true).
		MaxWidth(m.width).
		// Background(lipgloss.Color("#2c3e50")).
		BorderForeground(lipgloss.Color("#d67ab1")).
		Padding(0, 1)

	attachmentStyle := lipgloss.NewStyle().Faint(true)

	var content string
	breakpoints := " ,-"
	if m.attachmentName != "" {
		attachmentNameWrapped := ansi.Wordwrap(m.attachmentName, m.width, breakpoints)
		content = ansi.Wordwrap(m.prompt+"\n"+attachmentStyle.Render(attachmentNameWrapped), m.width, breakpoints)
	} else {
		content = ansi.Wordwrap(m.prompt, m.width, breakpoints)
	}

	return usermsgStyle.Render(content)
}

// AddUserMsgCmd returns a tea.Cmd that sends a UserMsgAddedMsg with the given prompt and attachmentName.
func AddUserMsgCmd(prompt string, attachmentName string) tea.Cmd {
	return func() tea.Msg {
		return UserMsgAddedMsg{
			UserMsg: NewUserMsg(prompt, attachmentName),
		}
	}
}

func NewUserMsg(prompt string, attachmentName string) UserMessage {
	return UserMessage{
		prompt:         prompt,
		attachmentName: attachmentName,
	}
}

/*
FIX: messages disappears on windowResize.
*/
