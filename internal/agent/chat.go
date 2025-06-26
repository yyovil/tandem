package agent

import (
	"context"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// NOTE: this is what the leftpane should look like after refactoring. this gots to be a bubble.
type Chat struct {
	ChatId   string
	Title    string
	viewport viewport.Model
	Channel  <-chan Message
	History  []Message
	agent    Agent
}

type ChatKeyMap struct {
	StopAgent,
	UpdateHistory key.Binding
}

var chatKeyMap = ChatKeyMap{
	StopAgent: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "stop agent"),
	),
}

func (c Chat) Init() tea.Cmd {}
func (c Chat) View() string  {}

func (c Chat) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	ctx, cancel := context.WithCancel(context.Background())
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case StreamCreated:
		c.Channel = msg.Stream
	case Message:
		c.History = append(c.History, msg)
		// runs agent as a side effect of appending the user message into the history.
		if msg.Type == UserMessageMsg {
			cmds = append(cmds, c.agent.Run(ctx, c.History))
		}

		cmds = append(cmds, c.Next())
		return c, tea.Batch(cmds...)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, chatKeyMap.StopAgent):
			cancel()
			// NOTE: I could put this inside a tea.Cmd to do the other things as well asynchronously ig but that's only if its required.
			return c, nil
		}
	}

	return c, tea.Batch(cmds...)
}

// pops them messages out of the Channel.
func (c Chat) Next() tea.Cmd {
	return func() tea.Msg {
		message, ok := <-c.Channel
		if !ok {
			// NOTE: its provider's responsibility to set the message type to ResponseCompletedMsg, representing end of the stream.
			return nil
		}
		return message
	}
}

func NewChat() Chat {
	return Chat{
		Channel: make(chan Message),
		History: []Message{},
	}
}
