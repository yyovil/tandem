package messages

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	// "github.com/yyovil/tui/internal/styles"
	"github.com/yyovil/tui/internal/utils"
)

type AgentMessage struct {
	StreamChan    chan tea.Msg
	Width, Height int
	Content       string
}

type AgentMessageAddedMsg struct {
	StreamChan chan tea.Msg
}

func (m *AgentMessage) Init() tea.Cmd {
	return nil
}

func (m *AgentMessage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ConcatenateChunkMsg:
		m.Content += string(msg)
		log.Printf("ConcatenateChunk: %s\n", msg)
		return m, ListenOnStreamChanCmd(m.StreamChan)
	case EndStream:
		log.Println("agentmessage: ending the stream.")
	}

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

	glammedResponse, err := glamour.Render(m.Content, "dark")
	if err != nil {
		log.Println("Error rendering response:", err)
		// TODO: provide a user feedback for this error
		return agentMessageStyle.Render(ansi.Wordwrap(m.Content, m.Width, utils.Breakpoints))
	}

	return glammedResponse
}

type ConcatenateChunkMsg string
type EndStream struct{}

func ListenOnStreamChanCmd(streamChan chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-streamChan
		if !ok {
			return EndStream{}
		}
		return msg
	}
}
