package styles

import (
	"github.com/charmbracelet/lipgloss"
)

/*
TODO:
1. here you define your glamour styles for the application.
2. some of your base styles could accommodate in here.
*/

var MessageStyle = lipgloss.
	NewStyle().
	Border(lipgloss.InnerHalfBlockBorder(), false).
	BorderLeft(true).
	Padding(0, 1)

