package messages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/yyovil/tui/internal/utils"
)

type UserMessage struct {
	prompt        string
	attachments   []string
	Width, Height int
}

// UserMessageAddedMsg is a message indicating a user message should be added to the leftpane viewport.
type UserMessageAddedMsg struct {
	UserMessage UserMessage
}

func (m *UserMessage) Init() tea.Cmd {
	return nil
}

func (m *UserMessage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *UserMessage) View() string {
	userMessageStyle := lipgloss.
		NewStyle().
		Border(lipgloss.InnerHalfBlockBorder(), false, false, false, true).
		MaxWidth(m.Width).
		// Background(lipgloss.Color("#2c3e50")).
		BorderForeground(lipgloss.Color("#ffafcc")).
		Padding(0, 1)

	attachmentStyle := lipgloss.NewStyle().Faint(true)
	var content string
	if len(m.attachments) > 0 {
		attachmentNameWrapped := ansi.Wordwrap(lipgloss.JoinVertical(lipgloss.Top, m.attachments...), m.Width, utils.Breakpoints)
		content = ansi.Wordwrap(m.prompt+"\n"+attachmentStyle.Render(attachmentNameWrapped), m.Width-2, utils.Breakpoints)
	} else {
		content = ansi.Wordwrap(m.prompt, m.Width-2, utils.Breakpoints)
	}

	return userMessageStyle.Render(content)
}

func AddUserMessageCmd(prompt string, attachments []string) tea.Cmd {
	return func() tea.Msg {
		return UserMessageAddedMsg{
			UserMessage: NewUserMessage(prompt, attachments),
		}
	}
}

func NewUserMessage(prompt string, attachments []string) UserMessage {
	return UserMessage{
		prompt:      prompt,
		attachments: attachments,
	}
}
