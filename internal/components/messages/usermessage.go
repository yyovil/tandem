package messages

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yyovil/tandem/internal/utils"
)

type UserMessage struct {
	prompt        string
	attachments   []string
	Width, Height int
}

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
		Width(m.Width).
		MaxWidth(m.Width).
		Background(lipgloss.Color("#2A1F23")).
		BorderForeground(lipgloss.Color("#ffafcc")).
		Padding(0, 1)

	attachmentStyle := lipgloss.NewStyle().Faint(true)
	var content strings.Builder
	if len(m.attachments) > 0 {
		attachmentNameWrapped := utils.Wordwrap(lipgloss.JoinVertical(lipgloss.Top, m.attachments...), m.Width)
		content.WriteString(utils.Wordwrap(m.prompt+"\n"+attachmentStyle.Render(attachmentNameWrapped), m.Width-2))
	} else {
		content.WriteString(utils.Wordwrap(m.prompt, m.Width-2))
	}

	return userMessageStyle.Render(content.String())
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
