package app

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yyovil/tandem/internal/agent"
	"github.com/yyovil/tandem/internal/bubbles"
)

type Status string

const (
	Requesting    Status = "requesting"
	Streaming     Status = "streaming"
	ToolCall      Status = "tool_call"
	ToolCompleted Status = "tool_completed"
	Idle          Status = "idle"
	Error         Status = "error"
)

/*
!TODO:
1. provide help for a very intuitive user experience.
*/
type App struct {
	width, height int
	chat          bubbles.Chat //in-memory chat history that is sent to the agent.
	Status        Status
	input         bubbles.Input
	Dialog        bubbles.Dialog
	layout        bubbles.SplitPane // this is the split pane layout that renders the chat and input bubbles.
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
	input := &bubbles.Input{}
	cmd := input.Init()
	return cmd
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height

		a.chat.Height = int(float32(msg.Height) * a.layout.HeightRatio)
		a.chat.Width = int(float32(msg.Width) * a.layout.WidthRatio)

		a.input.Width = a.chat.Width
		a.input.Height = msg.Height - a.chat.Height

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, appKeyMap.SendMessage):
			/*
				!TODO:
				1. within a session there could be only one active agent.
			*/

			input, cmd := a.input.Update(msg)
			a.input = input.(bubbles.Input)
			cmds = append(cmds, cmd)

			updateHistoryMsg := agent.Message{
				Type:  agent.UserMessageMsg,
				Files: []agent.Blob{},
				Part: agent.Part{
					Text: a.input.UserPrompt,
				},
				Role: agent.RoleUser,
			}
			chatModel, cmd := a.chat.Update(updateHistoryMsg)
			a.chat = chatModel.(bubbles.Chat)
			cmds = append(cmds, cmd)

			return a, tea.Batch(cmds...)

		case key.Matches(msg, appKeyMap.SelectModel):
			/*
				!TODO:
				1. open selectModelDialog
				2. update selectModelDialog
			*/
			return a, tea.Batch(cmds...)
		default:
			// NOTE: input should only get those key msgs when no dialog is open. because then its the dialog turn to receive them.
			model, cmd := a.input.Update(msg)
			a.input = model.(bubbles.Input)
			cmds = append(cmds, cmd)
		}
	}

	model, cmd := a.input.Update(msg)
	a.input = model.(bubbles.Input)
	cmds = append(cmds, cmd)

	model, cmd = a.chat.Update(msg)
	a.chat = model.(bubbles.Chat)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a *App) View() string {
	// NOTE: this way you can also control which pane to render depending on the terminal size and user preferences.
	a.layout.Leftpane = a.chat.View()
	a.layout.Bottom = a.input.View()
	a.layout.Status = string(a.Status)
	return a.layout.View()
}

func NewApp() *tea.Program {

	app := &App{
		// dialog:  bubbles.NewDialog(),
		Status: Idle,
		chat:   bubbles.NewChat(),
		input:  bubbles.NewInput(),
		layout: bubbles.SplitPane{
			WidthRatio:  70 / 100,
			HeightRatio: 80 / 100,
		},
	}

	return tea.NewProgram(app, tea.WithAltScreen(), tea.WithMouseCellMotion())
}
