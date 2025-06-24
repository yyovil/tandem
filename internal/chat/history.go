package chat

import (
	"github.com/charmbracelet/bubbletea"
)

type History struct {
	Channel  chan Message
	Messages []Message
}

func (h History) Next() tea.Cmd {
	return func() tea.Msg {
		content, ok := <-h.Channel
		if !ok {
			return Message{
				Type: ResponseCompletedMsg,
			}
		}
		return content
	}
}
