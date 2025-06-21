package layout

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	cbAnsi "github.com/charmbracelet/x/ansi"
	"github.com/muesli/ansi"
	"github.com/muesli/reflow/truncate"
	"github.com/yyovil/tandem/internal/utils"
)

const whitespace = " "

func getLines(s string) (lines []string, maxWidth int) {
	lines = strings.Split(s, "\n")
	for _, line := range lines {
		w := ansi.PrintableRuneWidth(line)
		if w > maxWidth {
			maxWidth = w
		}
	}

	return lines, maxWidth
}

// Composites the fg and the bg view.
func Composite(x, y int, fg, bg string) string {
	fgLines, fgWidth := getLines(fg)
	bglines, bgWidth := getLines(bg)
	fgHeight := len(fgLines)
	bgHeight := len(bglines)

	if fgWidth > bgWidth && fgHeight > bgHeight {
		return fg
	}

	x = utils.Clamp(x, 0, bgWidth-fgWidth)
	y = utils.Clamp(y, 0, bgHeight-fgHeight)

	var view strings.Builder

	for row, bgLine := range bglines {
		// ADHD: this is something mysterious. after this, compositing started working all of sudden. honestly, I don't know why. and the guy from whom i copied this copied from someone else in turn OG: github.com/yorukot you would get to see the PR mentioned in the opencode-ai/opencode.
		if row > 0 {
			view.WriteString("\n")
		}

		if row < y || row >= y+fgHeight {
			view.WriteString(bgLine)
			continue
		}

		col := 0
		if x > 0 {
			leftString := truncate.String(bgLine, uint(x))
			col = ansi.PrintableRuneWidth(leftString)
			view.WriteString(leftString)

			if col < x {
				view.WriteString(strings.Repeat(whitespace, x-col))
				col = x
			}
		}

		fgLine := fgLines[row-y]
		view.WriteString(fgLine)
		col += ansi.PrintableRuneWidth(fgLine)

		rightString := getRightString(bgLine, col)
		rightWidth := ansi.PrintableRuneWidth(rightString)

		if rightWidth <= bgWidth-col {
			view.WriteString(strings.Repeat(whitespace, bgWidth-col-rightWidth))
		}

		view.WriteString(rightString)
	}

	return view.String()
}

func getRightString(s string, width int) string {
	return cbAnsi.Cut(s, width, lipgloss.Width(s))
}
