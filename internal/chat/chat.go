package chat

import (
	"github.com/charmbracelet/bubbles/viewport"
)

type Chat struct {
	ChatId   string
	Title    string
	viewport viewport.Model
	History  History
}
