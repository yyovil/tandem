package messages

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
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
		m.Content += string(msg)
		return m, ListenOnStreamChanCmd(m.StreamChan)
	}

	return m, nil
}

func (m *AgentMessage) View() string {
	agentMessageStyle := lipgloss.
		NewStyle().
		Border(lipgloss.InnerHalfBlockBorder(), false, false, false, true).
		MaxWidth(m.Width).
		Padding(0, 1)

	glam, err := glamour.NewTermRenderer()
	if err != nil {
		log.Println("Error creating glamour renderer:", err)
		// TODO: provide a user feedback for this error
		return agentMessageStyle.Render(m.Content)
	}

	glammedResponse, err := glam.Render(m.Content)
	if err != nil {
		log.Println("Error rendering response:", err)
		// TODO: provide a user feedback for this error
		return agentMessageStyle.Render(m.Content)
	}

	return agentMessageStyle.Render(glammedResponse)
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
			log.Println("Received message from stream channel:", msg)
			return msg
		} else {
			return EndStream{}
		}
	}
}
