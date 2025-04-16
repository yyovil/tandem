package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"net/http"
	"os"
)

// this model reprs your entire state of the cli app.
type Model struct {
	message,
	sessionId,
	userId,
	/*
		TODO: user should be able to select the model by himself because you don't know how they be feeling some type of way. create a ENUM for models.
	*/
	model string
}

// Request model for running an agent
type RunRequest struct {
	sessionId,
	model,
	userId,
	message string
	stream bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			// we have to make the post req.
			runRequest := RunRequest{
				message:   m.message,
				model:     m.model,
				sessionId: m.sessionId,
				userId:    m.userId,
			}

			requestBody, err := json.Marshal(runRequest)
			if err != nil {
				// wouldn't it be nice to tell the reason why it failed during in the dev env?
				return m, tea.Quit
			}

			// format the endpoint str.
			resp, err := http.Post(os.Getenv("ENDPOINT"), "application/json", bytes.NewBuffer(requestBody))
			if err != nil {
				return m, tea.Quit
			}

			defer resp.Body.Close()
		}
	}
	return m, nil
}

func (m Model) View() string {}

// func (m Model) Stream(msg completionOutput) tea.Cmd{
// 	return func() {}
// }

func initialModel() tea.Model {
	return Model{
		message:   "",
		model:     "gemini-2.0-flash-lite",
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
	append()
}


/*
TODO:
	>  use errors pkgs to provide better error handling instead of trying to log to stdout as that is being occupied by the tui.
	>  use a logger to log the errors to a file.
*/
