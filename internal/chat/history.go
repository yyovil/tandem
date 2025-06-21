package chat

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/genai"
)

// NOTE: this is what the leftpane should look like after refactoring.
type History struct {
	Channel  chan genai.Content
	Width    int
	Messages []tea.Msg
}

func (h *History) Init() tea.Cmd {
	return nil
}

func (h *History) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Implement your update logic here
	switch msg := msg.(type) {
	case AddUserMsg:
		h.Messages = append(h.Messages, msg)
		return h, h.Next()

	case ResponseCompletedMsg:
		return h, nil
	}

	return h, nil
}

func (h *History) View() string {
	// Implement your view rendering here
	var view strings.Builder
	for _, message := range h.Messages {
		// TODO: style them messages, render them, join them vertically.
	}
	return view.String()
}

func NewHistory() *History {
	return &History{
		Width:    0,
		Messages: []tea.Msg{},
		Channel:  make(chan genai.Content),
	}
}

func (h *History) Next() tea.Cmd {
	return func() tea.Msg {
		content, ok := <-h.Channel
		if !ok {
			return ResponseCompletedMsg{}
		}
		return content
	}
}

/*
NOTE:
History contains the all the messages produced in a chat session.
so why not simply use it for rendering the view as well.

if you want to add parts to the message, dispatch a message and frwd it to the history cmp.
you want to add usermessage, just dispatch it towards the history cmp bro.

*/
