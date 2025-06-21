package messages

import (
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/yyovil/tandem/internal/utils"
)

type AgentMessage struct {
	Width, Height int
	Content       *strings.Builder
}
type RunStartedMsg struct{}
type RunResponseContentMsg RunResponse
type RunResponseCompletedMsg RunResponse

func (m *AgentMessage) Init() tea.Cmd {
	return nil
}

func (m *AgentMessage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *AgentMessage) View() string {
	agentMessageStyle := lipgloss.
		NewStyle().
		MaxWidth(m.Width).
		Border(lipgloss.InnerHalfBlockBorder(), false).
		BorderLeft(true).
		Background(lipgloss.Color("#cdb4db")).
		BorderForeground(lipgloss.Color("#cdb4db"))

	glammedResponse, err := glamour.Render(m.Content.String(), "dark")

	if err != nil {
		log.Println("Error rendering response:", err)
		// TODO: provide a user feedback for this error
		return agentMessageStyle.Render(utils.Wordwrap(m.Content.String(), m.Width))
	}

	return glammedResponse
}

func ListenOnStreamChanCmd(streamChan <-chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-streamChan
		if !ok {
			return RunResponseCompletedMsg{}
		}
		return msg
	}
}
