package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

// this model reprs your entire state of the cli app.
type Model struct {
	textarea textarea.Model
	message,
	sessionId,
	userId,
	streamingResponse,
	/*
		TODO: user should be able to select the model by himself because you don't know how they be feeling some type of way. create a ENUM for models.
	*/
	model string
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyTab:
			// TODO: make a request to the server to get the response stream.
			return m, m.GetCompletionStreamCmd()
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}
	}
	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// update the View to render the streaming response.
func (m Model) View() string {
	return fmt.Sprintf(
		"Chat with Sage.\n\n%s\n\n%s",
		m.textarea.View(),
		"ctrl+c: quit | tab: to send",
	) + "\n\n"
}

func initialModel() tea.Model {
	ta := textarea.New()
	ta.Placeholder = "Enter your prompt..."
	ta.Focus()

	return Model{
		textarea:  ta,
		message:   "",
		model:     ModelGeminiFlashLite,
		userId:    "slimeMaster",
		sessionId: "slimeMasterSession1",
	}
}

func main() {
	tuiLoop := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := tuiLoop.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

/*
TODO:
	>  use errors pkgs to provide better error handling instead of trying to log to stdout as that is being occupied by the tui.
	>  use a logger to log the errors to a file.
*/
