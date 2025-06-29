package bubbles

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/*
NOTE: This is a bubble that renders the split pane layout.

!TODO: viewports should be here and controlled from this bubble for the repective bubbles. bubbles should be designed in such a way that they can size to the viewport size.
*/
type SplitPane struct {
	WidthRatio  float32 // left pane: right pane ratio.
	HeightRatio float32 // left pane: bottom pane ratio.
	Leftpane    string
	Rightpane   string
	Bottom      string
	Status      string
}

func (s SplitPane) Init() tea.Cmd {
	return nil
}

// !TODO: we can use this to handle keybindings for the split pane, like resizing panes or switching focus.
func (s SplitPane) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return s, nil
}

func (s SplitPane) View() string {
	panes := lipgloss.JoinHorizontal(lipgloss.Left, s.Leftpane, s.Rightpane)
	layout := lipgloss.JoinVertical(lipgloss.Top, panes, s.Bottom)
	return layout
}
