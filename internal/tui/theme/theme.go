package theme

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme defines the interface for all UI themes in the application.
// All colors must be defined as lipgloss.AdaptiveColor to support
// both light and dark terminal backgrounds.
type Theme interface {
	// Diff view colors
	DiffAdded() lipgloss.AdaptiveColor
	DiffRemoved() lipgloss.AdaptiveColor
	DiffContext() lipgloss.AdaptiveColor
	DiffHunkHeader() lipgloss.AdaptiveColor
	DiffHighlightAdded() lipgloss.AdaptiveColor
	DiffHighlightRemoved() lipgloss.AdaptiveColor
	DiffAddedBg() lipgloss.AdaptiveColor
	DiffRemovedBg() lipgloss.AdaptiveColor
	DiffContextBg() lipgloss.AdaptiveColor
	DiffLineNumber() lipgloss.AdaptiveColor
	DiffAddedLineNumberBg() lipgloss.AdaptiveColor
	DiffRemovedLineNumberBg() lipgloss.AdaptiveColor

	// Base colors
	Primary() lipgloss.AdaptiveColor
	Secondary() lipgloss.AdaptiveColor
	Accent() lipgloss.AdaptiveColor

	// Status colors
	Error() lipgloss.AdaptiveColor
	Warning() lipgloss.AdaptiveColor
	Success() lipgloss.AdaptiveColor
	Info() lipgloss.AdaptiveColor

	// Text colors
	Text() lipgloss.AdaptiveColor
	TextMuted() lipgloss.AdaptiveColor
	TextEmphasized() lipgloss.AdaptiveColor

	// Background colors
	Background() lipgloss.AdaptiveColor
	BackgroundSecondary() lipgloss.AdaptiveColor
	BackgroundDarker() lipgloss.AdaptiveColor

	// Border colors
	BorderNormal() lipgloss.AdaptiveColor
	BorderFocused() lipgloss.AdaptiveColor
	BorderDim() lipgloss.AdaptiveColor

	// Markdown colors
	MarkdownText() lipgloss.AdaptiveColor
	MarkdownHeading() lipgloss.AdaptiveColor
	MarkdownLink() lipgloss.AdaptiveColor
	MarkdownLinkText() lipgloss.AdaptiveColor
	MarkdownCode() lipgloss.AdaptiveColor
	MarkdownBlockQuote() lipgloss.AdaptiveColor
	MarkdownEmph() lipgloss.AdaptiveColor
	MarkdownStrong() lipgloss.AdaptiveColor
	MarkdownHorizontalRule() lipgloss.AdaptiveColor
	MarkdownListItem() lipgloss.AdaptiveColor
	MarkdownListEnumeration() lipgloss.AdaptiveColor
	MarkdownImage() lipgloss.AdaptiveColor
	MarkdownImageText() lipgloss.AdaptiveColor
	MarkdownCodeBlock() lipgloss.AdaptiveColor

	// Syntax highlighting colors
	SyntaxComment() lipgloss.AdaptiveColor
	SyntaxKeyword() lipgloss.AdaptiveColor
	SyntaxFunction() lipgloss.AdaptiveColor
	SyntaxVariable() lipgloss.AdaptiveColor
	SyntaxString() lipgloss.AdaptiveColor
	SyntaxNumber() lipgloss.AdaptiveColor
	SyntaxType() lipgloss.AdaptiveColor
	SyntaxOperator() lipgloss.AdaptiveColor
	SyntaxPunctuation() lipgloss.AdaptiveColor
}

// BaseTheme provides a default implementation of the Theme interface
// that can be embedded in concrete theme implementations.
type BaseTheme struct {

	// Diff view colors
	DiffAddedColor               lipgloss.AdaptiveColor
	DiffRemovedColor             lipgloss.AdaptiveColor
	DiffContextColor             lipgloss.AdaptiveColor
	DiffHunkHeaderColor          lipgloss.AdaptiveColor
	DiffHighlightAddedColor      lipgloss.AdaptiveColor
	DiffHighlightRemovedColor    lipgloss.AdaptiveColor
	DiffAddedBgColor             lipgloss.AdaptiveColor
	DiffRemovedBgColor           lipgloss.AdaptiveColor
	DiffContextBgColor           lipgloss.AdaptiveColor
	DiffLineNumberColor          lipgloss.AdaptiveColor
	DiffAddedLineNumberBgColor   lipgloss.AdaptiveColor
	DiffRemovedLineNumberBgColor lipgloss.AdaptiveColor

	// Base colors
	PrimaryColor   lipgloss.AdaptiveColor
	SecondaryColor lipgloss.AdaptiveColor
	AccentColor    lipgloss.AdaptiveColor

	// Status colors
	ErrorColor   lipgloss.AdaptiveColor
	WarningColor lipgloss.AdaptiveColor
	SuccessColor lipgloss.AdaptiveColor
	InfoColor    lipgloss.AdaptiveColor

	// Text colors
	TextColor           lipgloss.AdaptiveColor
	TextMutedColor      lipgloss.AdaptiveColor
	TextEmphasizedColor lipgloss.AdaptiveColor

	// Background colors
	BackgroundColor          lipgloss.AdaptiveColor
	BackgroundSecondaryColor lipgloss.AdaptiveColor
	BackgroundDarkerColor    lipgloss.AdaptiveColor

	// Border colors
	BorderNormalColor  lipgloss.AdaptiveColor
	BorderFocusedColor lipgloss.AdaptiveColor
	BorderDimColor     lipgloss.AdaptiveColor

	// Markdown colors
	MarkdownTextColor            lipgloss.AdaptiveColor
	MarkdownHeadingColor         lipgloss.AdaptiveColor
	MarkdownLinkColor            lipgloss.AdaptiveColor
	MarkdownLinkTextColor        lipgloss.AdaptiveColor
	MarkdownCodeColor            lipgloss.AdaptiveColor
	MarkdownBlockQuoteColor      lipgloss.AdaptiveColor
	MarkdownEmphColor            lipgloss.AdaptiveColor
	MarkdownStrongColor          lipgloss.AdaptiveColor
	MarkdownHorizontalRuleColor  lipgloss.AdaptiveColor
	MarkdownListItemColor        lipgloss.AdaptiveColor
	MarkdownListEnumerationColor lipgloss.AdaptiveColor
	MarkdownImageColor           lipgloss.AdaptiveColor
	MarkdownImageTextColor       lipgloss.AdaptiveColor
	MarkdownCodeBlockColor       lipgloss.AdaptiveColor

	// Syntax highlighting colors
	SyntaxCommentColor     lipgloss.AdaptiveColor
	SyntaxKeywordColor     lipgloss.AdaptiveColor
	SyntaxFunctionColor    lipgloss.AdaptiveColor
	SyntaxVariableColor    lipgloss.AdaptiveColor
	SyntaxStringColor      lipgloss.AdaptiveColor
	SyntaxNumberColor      lipgloss.AdaptiveColor
	SyntaxTypeColor        lipgloss.AdaptiveColor
	SyntaxOperatorColor    lipgloss.AdaptiveColor
	SyntaxPunctuationColor lipgloss.AdaptiveColor
}

// Implement the Theme interface for BaseTheme
func (t *BaseTheme) Primary() lipgloss.AdaptiveColor   { return t.PrimaryColor }
func (t *BaseTheme) Secondary() lipgloss.AdaptiveColor { return t.SecondaryColor }
func (t *BaseTheme) Accent() lipgloss.AdaptiveColor    { return t.AccentColor }

func (t *BaseTheme) Error() lipgloss.AdaptiveColor   { return t.ErrorColor }
func (t *BaseTheme) Warning() lipgloss.AdaptiveColor { return t.WarningColor }
func (t *BaseTheme) Success() lipgloss.AdaptiveColor { return t.SuccessColor }
func (t *BaseTheme) Info() lipgloss.AdaptiveColor    { return t.InfoColor }

func (t *BaseTheme) Text() lipgloss.AdaptiveColor           { return t.TextColor }
func (t *BaseTheme) TextMuted() lipgloss.AdaptiveColor      { return t.TextMutedColor }
func (t *BaseTheme) TextEmphasized() lipgloss.AdaptiveColor { return t.TextEmphasizedColor }

func (t *BaseTheme) Background() lipgloss.AdaptiveColor          { return t.BackgroundColor }
func (t *BaseTheme) BackgroundSecondary() lipgloss.AdaptiveColor { return t.BackgroundSecondaryColor }
func (t *BaseTheme) BackgroundDarker() lipgloss.AdaptiveColor    { return t.BackgroundDarkerColor }

func (t *BaseTheme) BorderNormal() lipgloss.AdaptiveColor  { return t.BorderNormalColor }
func (t *BaseTheme) BorderFocused() lipgloss.AdaptiveColor { return t.BorderFocusedColor }
func (t *BaseTheme) BorderDim() lipgloss.AdaptiveColor     { return t.BorderDimColor }

func (t *BaseTheme) MarkdownText() lipgloss.AdaptiveColor       { return t.MarkdownTextColor }
func (t *BaseTheme) MarkdownHeading() lipgloss.AdaptiveColor    { return t.MarkdownHeadingColor }
func (t *BaseTheme) MarkdownLink() lipgloss.AdaptiveColor       { return t.MarkdownLinkColor }
func (t *BaseTheme) MarkdownLinkText() lipgloss.AdaptiveColor   { return t.MarkdownLinkTextColor }
func (t *BaseTheme) MarkdownCode() lipgloss.AdaptiveColor       { return t.MarkdownCodeColor }
func (t *BaseTheme) MarkdownBlockQuote() lipgloss.AdaptiveColor { return t.MarkdownBlockQuoteColor }
func (t *BaseTheme) MarkdownEmph() lipgloss.AdaptiveColor       { return t.MarkdownEmphColor }
func (t *BaseTheme) MarkdownStrong() lipgloss.AdaptiveColor     { return t.MarkdownStrongColor }
func (t *BaseTheme) MarkdownHorizontalRule() lipgloss.AdaptiveColor {
	return t.MarkdownHorizontalRuleColor
}
func (t *BaseTheme) MarkdownListItem() lipgloss.AdaptiveColor { return t.MarkdownListItemColor }
func (t *BaseTheme) MarkdownListEnumeration() lipgloss.AdaptiveColor {
	return t.MarkdownListEnumerationColor
}
func (t *BaseTheme) MarkdownImage() lipgloss.AdaptiveColor     { return t.MarkdownImageColor }
func (t *BaseTheme) MarkdownImageText() lipgloss.AdaptiveColor { return t.MarkdownImageTextColor }
func (t *BaseTheme) MarkdownCodeBlock() lipgloss.AdaptiveColor { return t.MarkdownCodeBlockColor }

func (t *BaseTheme) SyntaxComment() lipgloss.AdaptiveColor     { return t.SyntaxCommentColor }
func (t *BaseTheme) SyntaxKeyword() lipgloss.AdaptiveColor     { return t.SyntaxKeywordColor }
func (t *BaseTheme) SyntaxFunction() lipgloss.AdaptiveColor    { return t.SyntaxFunctionColor }
func (t *BaseTheme) SyntaxVariable() lipgloss.AdaptiveColor    { return t.SyntaxVariableColor }
func (t *BaseTheme) SyntaxString() lipgloss.AdaptiveColor      { return t.SyntaxStringColor }
func (t *BaseTheme) SyntaxNumber() lipgloss.AdaptiveColor      { return t.SyntaxNumberColor }
func (t *BaseTheme) SyntaxType() lipgloss.AdaptiveColor        { return t.SyntaxTypeColor }
func (t *BaseTheme) SyntaxOperator() lipgloss.AdaptiveColor    { return t.SyntaxOperatorColor }
func (t *BaseTheme) SyntaxPunctuation() lipgloss.AdaptiveColor { return t.SyntaxPunctuationColor }

func (t *BaseTheme) DiffAdded() lipgloss.AdaptiveColor            { return t.DiffAddedColor }
func (t *BaseTheme) DiffRemoved() lipgloss.AdaptiveColor          { return t.DiffRemovedColor }
func (t *BaseTheme) DiffContext() lipgloss.AdaptiveColor          { return t.DiffContextColor }
func (t *BaseTheme) DiffHunkHeader() lipgloss.AdaptiveColor       { return t.DiffHunkHeaderColor }
func (t *BaseTheme) DiffHighlightAdded() lipgloss.AdaptiveColor   { return t.DiffHighlightAddedColor }
func (t *BaseTheme) DiffHighlightRemoved() lipgloss.AdaptiveColor { return t.DiffHighlightRemovedColor }
func (t *BaseTheme) DiffAddedBg() lipgloss.AdaptiveColor          { return t.DiffAddedBgColor }
func (t *BaseTheme) DiffRemovedBg() lipgloss.AdaptiveColor        { return t.DiffRemovedBgColor }
func (t *BaseTheme) DiffContextBg() lipgloss.AdaptiveColor        { return t.DiffContextBgColor }
func (t *BaseTheme) DiffLineNumber() lipgloss.AdaptiveColor       { return t.DiffLineNumberColor }
func (t *BaseTheme) DiffAddedLineNumberBg() lipgloss.AdaptiveColor {
	return t.DiffAddedLineNumberBgColor
}
func (t *BaseTheme) DiffRemovedLineNumberBg() lipgloss.AdaptiveColor {
	return t.DiffRemovedLineNumberBgColor
}

// TODO: Customise the Theme for tandem later.
// CurrentTheme creates a new instance of the Tandem theme.
func CurrentTheme() Theme {
	// Tandem color palette
	// Dark mode colors
	darkBackground := "#000000"
	darkCurrentLine := "#252525"
	darkSelection := "#303030"
	darkForeground := "#ffffff"
	darkComment := "#6a6a6a"
	darkPrimary := "#ff26b3"   // Primary orange/gold
	darkSecondary := "#1e7cff" // Secondary blue
	darkAccent := "#9d7cd8"    // Accent purple
	darkRed := "#ec1532"       // Error red
	darkOrange := "#f5a742"    // Warning orange
	darkGreen := "#59dc38"     // Success green
	darkCyan := "#56b6c2"      // Info cyan
	darkYellow := "#e5c07b"    // Emphasized text
	darkBorder := "#4b4c5c"    // Border color

	// Light mode colors
	lightBackground := "#f8f8f8"
	lightCurrentLine := "#f0f0f0"
	lightSelection := "#e5e5e6"
	lightForeground := "#2a2a2a"
	lightComment := "#8a8a8a"
	lightPrimary := "#3b7dd8"   // Primary blue
	lightSecondary := "#7b5bb6" // Secondary purple
	lightAccent := "#d68c27"    // Accent orange/gold
	lightRed := "#d1383d"       // Error red
	lightOrange := "#d68c27"    // Warning orange
	lightGreen := "#3d9a57"     // Success green
	lightCyan := "#318795"      // Info cyan
	lightYellow := "#b0851f"    // Emphasized text
	lightBorder := "#d3d3d3"    // Border color

	theme := &BaseTheme{}

	// Base colors
	theme.PrimaryColor = lipgloss.AdaptiveColor{
		Dark:  darkPrimary,
		Light: lightPrimary,
	}
	theme.SecondaryColor = lipgloss.AdaptiveColor{
		Dark:  darkSecondary,
		Light: lightSecondary,
	}
	theme.AccentColor = lipgloss.AdaptiveColor{
		Dark:  darkAccent,
		Light: lightAccent,
	}

	// Status colors
	theme.ErrorColor = lipgloss.AdaptiveColor{
		Dark:  darkRed,
		Light: lightRed,
	}
	theme.WarningColor = lipgloss.AdaptiveColor{
		Dark:  darkOrange,
		Light: lightOrange,
	}
	theme.SuccessColor = lipgloss.AdaptiveColor{
		Dark:  darkGreen,
		Light: lightGreen,
	}
	theme.InfoColor = lipgloss.AdaptiveColor{
		Dark:  darkCyan,
		Light: lightCyan,
	}

	// Text colors
	theme.TextColor = lipgloss.AdaptiveColor{
		Dark:  darkForeground,
		Light: lightForeground,
	}
	theme.TextMutedColor = lipgloss.AdaptiveColor{
		Dark:  darkComment,
		Light: lightComment,
	}
	theme.TextEmphasizedColor = lipgloss.AdaptiveColor{
		Dark:  darkYellow,
		Light: lightYellow,
	}

	// Background colors
	theme.BackgroundColor = lipgloss.AdaptiveColor{
		Dark:  darkBackground,
		Light: lightBackground,
	}
	theme.BackgroundSecondaryColor = lipgloss.AdaptiveColor{
		Dark:  darkCurrentLine,
		Light: lightCurrentLine,
	}
	theme.BackgroundDarkerColor = lipgloss.AdaptiveColor{
		Dark:  "#121212", // Slightly darker than background
		Light: "#ffffff", // Slightly lighter than background
	}

	// Border colors
	theme.BorderNormalColor = lipgloss.AdaptiveColor{
		Dark:  darkBorder,
		Light: lightBorder,
	}
	theme.BorderFocusedColor = lipgloss.AdaptiveColor{
		Dark:  darkPrimary,
		Light: lightPrimary,
	}
	theme.BorderDimColor = lipgloss.AdaptiveColor{
		Dark:  darkSelection,
		Light: lightSelection,
	}

	// Markdown colors
	theme.MarkdownTextColor = lipgloss.AdaptiveColor{
		Dark:  darkForeground,
		Light: lightForeground,
	}
	theme.MarkdownHeadingColor = lipgloss.AdaptiveColor{
		Dark:  darkSecondary,
		Light: lightSecondary,
	}
	theme.MarkdownLinkColor = lipgloss.AdaptiveColor{
		Dark:  darkPrimary,
		Light: lightPrimary,
	}
	theme.MarkdownLinkTextColor = lipgloss.AdaptiveColor{
		Dark:  darkCyan,
		Light: lightCyan,
	}
	theme.MarkdownCodeColor = lipgloss.AdaptiveColor{
		Dark:  darkGreen,
		Light: lightGreen,
	}
	theme.MarkdownBlockQuoteColor = lipgloss.AdaptiveColor{
		Dark:  darkYellow,
		Light: lightYellow,
	}
	theme.MarkdownEmphColor = lipgloss.AdaptiveColor{
		Dark:  darkYellow,
		Light: lightYellow,
	}
	theme.MarkdownStrongColor = lipgloss.AdaptiveColor{
		Dark:  darkAccent,
		Light: lightAccent,
	}
	theme.MarkdownHorizontalRuleColor = lipgloss.AdaptiveColor{
		Dark:  darkComment,
		Light: lightComment,
	}
	theme.MarkdownListItemColor = lipgloss.AdaptiveColor{
		Dark:  darkPrimary,
		Light: lightPrimary,
	}
	theme.MarkdownListEnumerationColor = lipgloss.AdaptiveColor{
		Dark:  darkCyan,
		Light: lightCyan,
	}
	theme.MarkdownImageColor = lipgloss.AdaptiveColor{
		Dark:  darkPrimary,
		Light: lightPrimary,
	}
	theme.MarkdownImageTextColor = lipgloss.AdaptiveColor{
		Dark:  darkCyan,
		Light: lightCyan,
	}
	theme.MarkdownCodeBlockColor = lipgloss.AdaptiveColor{
		Dark:  darkForeground,
		Light: lightForeground,
	}

	// Syntax highlighting colors
	theme.SyntaxCommentColor = lipgloss.AdaptiveColor{
		Dark:  darkComment,
		Light: lightComment,
	}
	theme.SyntaxKeywordColor = lipgloss.AdaptiveColor{
		Dark:  darkSecondary,
		Light: lightSecondary,
	}
	theme.SyntaxFunctionColor = lipgloss.AdaptiveColor{
		Dark:  darkPrimary,
		Light: lightPrimary,
	}
	theme.SyntaxVariableColor = lipgloss.AdaptiveColor{
		Dark:  darkRed,
		Light: lightRed,
	}
	theme.SyntaxStringColor = lipgloss.AdaptiveColor{
		Dark:  darkGreen,
		Light: lightGreen,
	}
	theme.SyntaxNumberColor = lipgloss.AdaptiveColor{
		Dark:  darkAccent,
		Light: lightAccent,
	}
	theme.SyntaxTypeColor = lipgloss.AdaptiveColor{
		Dark:  darkYellow,
		Light: lightYellow,
	}
	theme.SyntaxOperatorColor = lipgloss.AdaptiveColor{
		Dark:  darkCyan,
		Light: lightCyan,
	}
	theme.SyntaxPunctuationColor = lipgloss.AdaptiveColor{
		Dark:  darkForeground,
		Light: lightForeground,
	}

	// Diff view colors
	theme.DiffAddedColor = lipgloss.AdaptiveColor{
		Dark:  "#478247",
		Light: "#2E7D32",
	}
	theme.DiffRemovedColor = lipgloss.AdaptiveColor{
		Dark:  "#7C4444",
		Light: "#C62828",
	}
	theme.DiffContextColor = lipgloss.AdaptiveColor{
		Dark:  "#a0a0a0",
		Light: "#757575",
	}
	theme.DiffHunkHeaderColor = lipgloss.AdaptiveColor{
		Dark:  "#a0a0a0",
		Light: "#757575",
	}
	theme.DiffHighlightAddedColor = lipgloss.AdaptiveColor{
		Dark:  "#DAFADA",
		Light: "#A5D6A7",
	}
	theme.DiffHighlightRemovedColor = lipgloss.AdaptiveColor{
		Dark:  "#FADADD",
		Light: "#EF9A9A",
	}
	theme.DiffAddedBgColor = lipgloss.AdaptiveColor{
		Dark:  "#303A30",
		Light: "#E8F5E9",
	}
	theme.DiffRemovedBgColor = lipgloss.AdaptiveColor{
		Dark:  "#3A3030",
		Light: "#FFEBEE",
	}
	theme.DiffContextBgColor = lipgloss.AdaptiveColor{
		Dark:  darkBackground,
		Light: lightBackground,
	}
	theme.DiffLineNumberColor = lipgloss.AdaptiveColor{
		Dark:  "#888888",
		Light: "#9E9E9E",
	}
	theme.DiffAddedLineNumberBgColor = lipgloss.AdaptiveColor{
		Dark:  "#293229",
		Light: "#C8E6C9",
	}
	theme.DiffRemovedLineNumberBgColor = lipgloss.AdaptiveColor{
		Dark:  "#332929",
		Light: "#FFCDD2",
	}

	return theme
}
