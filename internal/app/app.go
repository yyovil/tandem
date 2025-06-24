package app

import (
	"context"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yyovil/tandem/internal/agent"
	"github.com/yyovil/tandem/internal/chat"
	"github.com/yyovil/tandem/internal/components"
)

type Status string

const (
	Requesting    Status = "requesting"
	Streaming     Status = "streaming"
	ToolCall      Status = "tool_call"
	ToolCompleted Status = "tool_completed"
	Idle          Status = "idle"
)

type App struct {
	chat    chat.Chat //in-memory chat history that is sent to the agent.
	agent   agent.Agent
	status  Status
	input   components.Input
	dialog  components.Dialog
	context context.Context
	cancel  context.CancelFunc
}

type AppKeyMap struct {
	SelectModel,
	SendMessage key.Binding
}

var appKeyMap = AppKeyMap{
	SendMessage: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "send message"),
	),
	SelectModel: key.NewBinding(
		key.WithKeys("m", "ctrl+m"),
		key.WithHelp("m", "select model"),
	),
}

func (a *App) Init() tea.Cmd {
	a.context, a.cancel = context.WithCancel(context.Background())
	input := &components.Input{}
	input.Init()
	return nil
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: listen for messages from the agent.
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, appKeyMap.SendMessage):
			/*
				TODO:
				1. update the chat history.
				2. within a session there could be only one active agent.
			*/
			input, cmd := a.input.Update(msg)
			a.input = input.(components.Input)
			cmds = append(cmds, cmd)

			chatModel, cmd := a.chat.Update(msg)
			a.chat = chatModel.(chat.Chat)
			cmds = append(cmds, cmd)

			return a, tea.Batch(cmds...)

		case key.Matches(msg, appKeyMap.SelectModel):
			/*
				TODO:
				1. open selectModelDialog
				2. update selectModelDialog
			*/
			return a, tea.Batch(cmds...)
		}

		// TODO: forward keystrokes to the input component
	}

	return a, nil
}

func (a *App) View() string {
	return "Tandem App"
}

func NewApp() *tea.Program {

	// TODO: load the user settings here. also define fallback defaults.

	agent, err := agent.NewAgent()
	if err != nil {
		// TODO: produce better response for the user explaining the difficulty here.
		return nil
	}

	app := &App{
		status: Idle,
		agent:  agent,
		input:  components.NewInput(),
	}

	return tea.NewProgram(app, tea.WithAltScreen(), tea.WithMouseCellMotion())
}
