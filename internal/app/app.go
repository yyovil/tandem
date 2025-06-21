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
	Requesting    Status = "Requesting"
	Streaming     Status = "Streaming"
	ToolCall      Status = "Tool call"
	ToolCompleted Status = "Tool completed"
	Idle          Status = "Idle"
)

type App struct {
	history chat.History //in-memory chat history that is sent to the agent.
	agent   agent.Agent
	status  Status
	input   components.Input
	dialog  components.Dialog
	context context.Context
	cancel  context.CancelFunc
}

type AppKeyMap struct {
	SelectModel,
	RunAgent key.Binding
}

var appKeyMap = AppKeyMap{
	RunAgent: key.NewBinding(
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

	if agent, err := agent.NewAgent(a.setting); err != nil {
		// TODO: handle the error gracefully, maybe show a message to the user.
	} else {
		a.agent = agent
	}

	a.history.Init()

	input := &components.Input{}
	input.Init()

	return nil
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: listen for messages from the agent.
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, appKeyMap.RunAgent):
			/*
				TODO:
				1. update the chat history.
				2. clear the input textarea.
				3. within a session there could be only one active agent.
			*/
			return a, a.agent.Run(a.context, a.history.Channel)
		case key.Matches(msg, appKeyMap.SelectModel):
			/*
				TODO:
				1. open the dialog.

			*/

			_, cmd = a.setting.Update(msg)
			cmds = append(cmds, cmd)
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

	agent, err := agent.NewAgent()
	if err != nil {
		// TODO: produce better response for the user explaining the difficulty here.
		return nil
	}

	app := &App{
		status:  Idle,
		agent:   agent,
		input:   components.NewInput(),
		history: *chat.NewHistory(),
	}

	return tea.NewProgram(app, tea.WithAltScreen(), tea.WithMouseCellMotion())
}
