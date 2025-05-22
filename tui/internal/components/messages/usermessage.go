package messages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type UserMessage struct {
	prompt         string
	attachmentName string
	Width, Height  int
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
		BorderForeground(lipgloss.Color("#d67ab1")).
		Padding(0, 1)

	attachmentStyle := lipgloss.NewStyle().Faint(true)

	var content string
	breakpoints := " ,-"
	if m.attachmentName != "" {
		attachmentNameWrapped := ansi.Wordwrap(m.attachmentName, m.Width, breakpoints)
		content = ansi.Wordwrap(m.prompt+"\n"+attachmentStyle.Render(attachmentNameWrapped), m.Width-1, breakpoints)
	} else {
		content = ansi.Wordwrap(m.prompt, m.Width, breakpoints)
	}

	return userMessageStyle.Render(content)
}

// AddUserMessageCmd returns a tea.Cmd that sends a UserMsgAddedMsg with the given prompt and attachmentName.
func AddUserMessageCmd(prompt string, attachmentName string) tea.Cmd {
	return func() tea.Msg {
		return UserMessageAddedMsg{
			UserMessage: NewUserMessage(prompt, attachmentName),
		}
	}
}

func NewUserMessage(prompt string, attachmentName string) UserMessage {
	return UserMessage{
		prompt:         prompt,
		attachmentName: attachmentName,
	}
}
