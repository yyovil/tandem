package bubbles

import (
	"context"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/yyovil/tandem/internal/agent"
	"github.com/yyovil/tandem/internal/styles"
	"github.com/yyovil/tandem/internal/utils"
	"strings"
)

// NOTE: this is what the leftpane should look like after refactoring. this gots to be a bubble.
type Chat struct {
	Width, Height int
	ChatId        string
	Title         string
	Channel       <-chan agent.Message
	History       []agent.Message
	agent         agent.Agent
}

type ChatKeyMap struct {
	RunAgent,
	StopAgent key.Binding
}

var chatKeyMap = ChatKeyMap{
	StopAgent: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "stop agent"),
	),
}

func (c Chat) Init() tea.Cmd {
	return nil
}

func (c Chat) View() string {
	var chatHistory strings.Builder

	userMessageStyle := styles.MessageStyle.
		Width(c.Width).
		MaxWidth(c.Width).
		Background(lipgloss.Color("#2A1F23")).
		BorderForeground(lipgloss.Color("#ffafcc"))

	filesStyle := lipgloss.NewStyle().Faint(true)

	// toolCallStyle := styles.MessageStyle.
	// 	Width(c.Width).
	// 	MaxWidth(c.Width).
	// 	Background(lipgloss.Color("#0F160D")).
	// 	BorderForeground(lipgloss.Color("#b2ff9e"))

	agentMessageStyle := styles.MessageStyle.
		Width(c.Width).
		MaxWidth(c.Width).
		Background(lipgloss.Color("#cdb4db")).
		BorderForeground(lipgloss.Color("#cdb4db"))

	for _, msg := range c.History {
		var content strings.Builder
		switch msg.Type {
		case agent.UserMessageMsg:
			if len(msg.Files) > 0 {
				attachments := make([]string, 0, len(msg.Files))
				for _, file := range msg.Files {
					attachments = append(attachments, file.Name)
				}
				filenameWrapped := utils.Wordwrap(lipgloss.JoinVertical(lipgloss.Top, attachments...), c.Width)
				content.WriteString(utils.Wordwrap(msg.Part.Text+"\n"+filesStyle.Render(filenameWrapped), c.Width-2))
			} else {
				content.WriteString(utils.Wordwrap(msg.Part.Text, c.Width-2))
			}

			chatHistory.WriteString(userMessageStyle.Render(content.String()))

		case agent.ResponseMsg:
			text := msg.Part.Text
			glammedResponse, err := glamour.Render(text, "dark")
			if err != nil {
				agentMessageStyle.Render(text)
			}

			chatHistory.WriteString(glammedResponse)

		case agent.ToolCallMsg:
		case agent.ToolResponseMsg:
		case agent.ToolCallErrorMsg:
		case agent.ResponseCompletedMsg:
		}
	}

	return chatHistory.String()
}

func (c Chat) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	ctx, cancel := context.WithCancel(context.Background())
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case agent.StreamCreated:
		c.Channel = msg.Stream
	case agent.Message:
		c.History = append(c.History, msg)
		// runs agent as a side effect of appending the user message into the history.
		if msg.Type == agent.UserMessageMsg {
			cmds = append(cmds, c.agent.Run(ctx, c.History))
		}

		cmds = append(cmds, c.Next())
		return c, tea.Batch(cmds...)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, chatKeyMap.StopAgent):
			cancel()
			return c, nil
		}
	}

	return c, tea.Batch(cmds...)
}

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
		Channel: make(chan agent.Message),
		History: []agent.Message{},
	}
}
