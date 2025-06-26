package app

import (
	"context"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yyovil/tandem/internal/agent"
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
	chat    agent.Chat //in-memory chat history that is sent to the agent.
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
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, appKeyMap.SendMessage):
			/*
				TODO:
				1. update the chat history.
				2. within a session there could be only one active agent.
			*/

			// updates the user prompt, collects all the attachments and clears the textarea.
			input, cmd := a.input.Update(msg)
			a.input = input.(components.Input)
			cmds = append(cmds, cmd)

			// update the chat history.
			updateHistoryMsg := agent.Message{
				Type:  agent.UserMessageMsg,
				Files: []agent.Blob{},
				Part: agent.Part{
					Text: a.input.UserPrompt,
				},
				Role: agent.RoleUser,
			}

			chatModel, cmd := a.chat.Update(updateHistoryMsg)
			a.chat = chatModel.(agent.Chat)
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

		// TODO: forward keystrokes to the input bubble.
		// TODO: forward rest of the messages to the chat bubble.
	}

	return a, nil
}

func (a *App) View() string {
	return "Tandem App"
}

func NewApp() *tea.Program {

	ctx, cancel := context.WithCancel(context.Background())
	app := &App{
		dialog:  components.NewDialog(),
		status:  Idle,
		chat:    agent.NewChat(),
		input:   components.NewInput(),
		context: ctx,
		cancel:  cancel,
	}

	return tea.NewProgram(app, tea.WithAltScreen(), tea.WithMouseCellMotion())
}
