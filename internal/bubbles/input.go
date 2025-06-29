package bubbles

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Input struct {
	Width, Height int
	UserPrompt    string
	textarea      textarea.Model
}

type InputKeyMap struct {
	Send key.Binding
}

var inputKeyMap = InputKeyMap{
	Send: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "send message"),
	),
}

func (i Input) Init() tea.Cmd {
	i.textarea.Placeholder = "Assign tasks to AI Agents here..."
	i.textarea.ShowLineNumbers = false
	return tea.Batch(
		textarea.Blink,
		i.textarea.Focus(),
	)
}

func (i Input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, inputKeyMap.Send):
			if i.textarea.Value() == "" {
				return i, nil
			}

			i.UserPrompt = i.textarea.Value()
			i.textarea.Reset()
			return i, tea.Batch(cmds...)
		}

	default:
		cmd = i.textarea.Focus()
		cmds = append(cmds, cmd)
		i.textarea, cmd = i.textarea.Update(msg)
		cmds = append(cmds, cmd)
		return i, tea.Batch(cmds...)
	}
	return i, tea.Batch(cmds...)
}

func (i Input) View() string {

	inputStyle := lipgloss.
		NewStyle().
		Width(i.Width-2).
		MaxWidth(i.Width).
		Height(4).
		MaxHeight(6).
		Border(lipgloss.NormalBorder(), true)

	return inputStyle.Render(i.textarea.View())
}

func NewInput() Input {
	return Input{
		Width:      0,
		Height:     0,
		UserPrompt: "",
		textarea:   textarea.New(),
	}
}
