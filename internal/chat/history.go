package chat

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// NOTE: this is what the leftpane should look like after refactoring.
type History struct {
	Channel  chan Message
	Width    int
	Messages []Message
}

func (h *History) Init() tea.Cmd {
	return nil
}

func (h *History) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case Message:
		switch msg.Event {
		case ResponseCompleted:
		case ToolCall:
			//TODO: handle all of them messages here eventwise.
		}

	}

	return h, nil
}

func (h *History) View() string {
	var view strings.Builder
	for _, message := range h.Messages {
		// TODO: style them messages eventwise, render them, join them vertically.
	}
	return view.String()
}

func NewHistory() *History {
	return &History{
		Width:    0,
		Messages: []Message{},
		Channel:  make(chan Message),
	}
}

func (h History) Next() tea.Cmd {
	return func() tea.Msg {
		content, ok := <-h.Channel
		if !ok {
			return Message{
				Event: ResponseCompleted,
			}
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
