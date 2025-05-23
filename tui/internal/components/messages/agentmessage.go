package messages

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	glam "github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
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
	// take care of the concatenation msg here.
	switch msg := msg.(type) {
	case ConcatenateChunkMsg:
		// log.Println("Concatenating chunk: ", msg)
		m.Content += string(msg)
		return m, ListenOnStreamChanCmd(m.StreamChan)
	}

	return m, nil
}

func (m *AgentMessage) View() string {
	agentMessageStyle := lipgloss.
		NewStyle().
		MaxWidth(m.Width).
		Padding(0, 1)

	glammedResponse, err := glam.Render(m.Content, "dark")

	if err != nil {
		log.Println("Error rendering response:", err)
		// TODO: provide a user feedback for this error
		return agentMessageStyle.Render(ansi.Wordwrap(m.Content, m.Width, utils.Breakpoints))
	}

	return agentMessageStyle.Render(ansi.Wordwrap(glammedResponse, m.Width, utils.Breakpoints))
}

// give this cmd when we have to concatenate a new chunk to the last agent message.
type ConcatenateChunkMsg string
type EndStream struct{}

func NewAgentMessage(completion string) AgentMessage {
	return AgentMessage{
		Content: completion,
	}
}

func ListenOnStreamChanCmd(streamChan <-chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		if msg, ok := <-streamChan; ok {
			return msg
		} else {
			return EndStream{}
		}
	}
}
