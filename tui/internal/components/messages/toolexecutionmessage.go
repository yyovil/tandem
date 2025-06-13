package messages

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yyovil/tui/internal/utils"
)

type ToolExecutionMessage struct {
	Width          int
	ToolCallName   string
	ToolCallResult string
	Content        string
	Event          RunEvent
}
type ToolCallStartedMsg RunResponse
type ToolCallCompletedMsg RunResponse

func (teMsg *ToolExecutionMessage) Init() tea.Cmd {
	return nil
}

func (teMsg *ToolExecutionMessage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return teMsg, nil
}

func (teMsg *ToolExecutionMessage) View() string {
	toolExecutionStyle := lipgloss.
		NewStyle().
		Width(teMsg.Width).
		MaxWidth(teMsg.Width).
		Border(lipgloss.InnerHalfBlockBorder(), false).
		BorderLeft(true).
		Background(lipgloss.Color("#0F160D")).
		BorderForeground(lipgloss.Color("#b2ff9e")).
		Padding(0, 1)

	var content strings.Builder

	content.WriteString(toolExecutionStyle.Bold(true).Render(utils.Wordwrap("Tool call: "+teMsg.ToolCallName, teMsg.Width-2) + "\n"))
	// truncate this to some 15-20 characters.
	content.WriteString(toolExecutionStyle.Render(utils.Wordwrap(teMsg.ToolCallResult, teMsg.Width-2) + "\n"))
	content.WriteString(toolExecutionStyle.Bold(true).Render(utils.Wordwrap(teMsg.Content, teMsg.Width-2)))

	return content.String()
}
