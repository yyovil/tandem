package dialog

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yyovil/tandem/internal/session"
	"github.com/yyovil/tandem/internal/tui/layout"
	"github.com/yyovil/tandem/internal/tui/styles"
	"github.com/yyovil/tandem/internal/tui/theme"
	"github.com/yyovil/tandem/internal/utils"
)

// SessionSelectedMsg is sent when a session is selected
type SessionSelectedMsg struct {
	Session session.Session
}

// SessionCreatedMsg is sent when a session is created
type SessionCreatedMsg struct {
	Session session.Session
}

// SessionUpdatedMsg is sent when a session is updated
type SessionUpdatedMsg struct {
	Session session.Session
}

// SessionDeletedMsg is sent when a session is deleted
type SessionDeletedMsg struct {
	SessionID string
}

// CloseSessionDialogMsg is sent when the session dialog is closed
type CloseSessionDialogMsg struct{}

// SessionDialog interface for the session switching dialog
type SessionDialog interface {
	tea.Model
	layout.Bindings
	SetSessions(sessions []session.Session)
	SetSelectedSession(sessionID string)
	SetSessionService(service session.Service)
}

type dialogMode int

const (
	modeList dialogMode = iota
	modeEdit
	modeCreate
	modeConfirmDelete
)

type sessionDialogCmp struct {
	sessions          []session.Session
	selectedIdx       int
	width             int
	height            int
	selectedSessionID string
	sessionService    session.Service
	
	// Input and mode management
	mode              dialogMode
	textInput         textinput.Model
	originalTitle     string
}

type sessionKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Enter    key.Binding
	Escape   key.Binding
	J        key.Binding
	K        key.Binding
	New      key.Binding
	Edit     key.Binding
	Delete   key.Binding
	Confirm  key.Binding
}

var sessionKeys = sessionKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "previous session"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "next session"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select session"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "close/cancel"),
	),
	J: key.NewBinding(
		key.WithKeys("j"),
		key.WithHelp("j", "next session"),
	),
	K: key.NewBinding(
		key.WithKeys("k"),
		key.WithHelp("k", "previous session"),
	),
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new session"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit session"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete session"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "confirm"),
	),
}

func (s *sessionDialogCmp) Init() tea.Cmd {
	s.setupTextInput()
	return nil
}

func (s *sessionDialogCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle input modes first
		switch s.mode {
		case modeEdit, modeCreate:
			return s.handleInputMode(msg)
		case modeConfirmDelete:
			return s.handleDeleteConfirmation(msg)
		default:
			return s.handleListMode(msg)
		}
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		s.textInput.Width = max(30, min(s.width-20, 50))
	}
	return s, nil
}

func (s *sessionDialogCmp) handleListMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, sessionKeys.Up) || key.Matches(msg, sessionKeys.K):
		if s.selectedIdx > 0 {
			s.selectedIdx--
		}
		return s, nil
	case key.Matches(msg, sessionKeys.Down) || key.Matches(msg, sessionKeys.J):
		if s.selectedIdx < len(s.sessions)-1 {
			s.selectedIdx++
		}
		return s, nil
	case key.Matches(msg, sessionKeys.Enter):
		if len(s.sessions) > 0 {
			return s, utils.CmdHandler(SessionSelectedMsg{
				Session: s.sessions[s.selectedIdx],
			})
		}
		return s, nil
	case key.Matches(msg, sessionKeys.New):
		return s.enterCreateMode(), nil
	case key.Matches(msg, sessionKeys.Edit):
		if len(s.sessions) > 0 {
			return s.enterEditMode(), nil
		}
		return s, nil
	case key.Matches(msg, sessionKeys.Delete):
		if len(s.sessions) > 0 {
			s.mode = modeConfirmDelete
			return s, nil
		}
		return s, nil
	case key.Matches(msg, sessionKeys.Escape):
		return s, utils.CmdHandler(CloseSessionDialogMsg{})
	}
	return s, nil
}

func (s *sessionDialogCmp) handleInputMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, sessionKeys.Enter):
		return s.submitInput()
	case key.Matches(msg, sessionKeys.Escape):
		s.mode = modeList
		s.textInput.SetValue("")
		return s, nil
	default:
		var cmd tea.Cmd
		s.textInput, cmd = s.textInput.Update(msg)
		return s, cmd
	}
}

func (s *sessionDialogCmp) handleDeleteConfirmation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, sessionKeys.Confirm):
		return s.deleteSession()
	case key.Matches(msg, sessionKeys.Escape):
		s.mode = modeList
		return s, nil
	}
	return s, nil
}

func (s *sessionDialogCmp) enterCreateMode() *sessionDialogCmp {
	s.mode = modeCreate
	s.textInput.SetValue("")
	s.textInput.Focus()
	return s
}

func (s *sessionDialogCmp) enterEditMode() *sessionDialogCmp {
	if len(s.sessions) > 0 {
		s.mode = modeEdit
		selectedSession := s.sessions[s.selectedIdx]
		s.originalTitle = selectedSession.Title
		s.textInput.SetValue(selectedSession.Title)
		s.textInput.Focus()
	}
	return s
}

func (s *sessionDialogCmp) submitInput() (tea.Model, tea.Cmd) {
	title := strings.TrimSpace(s.textInput.Value())
	if title == "" {
		return s, nil
	}

	if s.sessionService == nil {
		return s, utils.ReportError(fmt.Errorf("session service not available"))
	}

	s.textInput.Blur()
	
	switch s.mode {
	case modeCreate:
		return s.createSession(title)
	case modeEdit:
		return s.updateSession(title)
	}
	
	return s, nil
}

func (s *sessionDialogCmp) createSession(title string) (tea.Model, tea.Cmd) {
	s.mode = modeList
	s.textInput.SetValue("")
	
	return s, func() tea.Msg {
		ctx := context.Background()
		session, err := s.sessionService.Create(ctx, title)
		if err != nil {
			return utils.ReportError(err)
		}
		return SessionCreatedMsg{Session: session}
	}
}

func (s *sessionDialogCmp) updateSession(title string) (tea.Model, tea.Cmd) {
	s.mode = modeList
	s.textInput.SetValue("")
	
	if len(s.sessions) == 0 {
		return s, nil
	}
	
	selectedSession := s.sessions[s.selectedIdx]
	selectedSession.Title = title
	
	return s, func() tea.Msg {
		ctx := context.Background()
		session, err := s.sessionService.Save(ctx, selectedSession)
		if err != nil {
			return utils.ReportError(err)
		}
		return SessionUpdatedMsg{Session: session}
	}
}

func (s *sessionDialogCmp) deleteSession() (tea.Model, tea.Cmd) {
	s.mode = modeList
	
	if len(s.sessions) == 0 {
		return s, nil
	}
	
	selectedSession := s.sessions[s.selectedIdx]
	
	return s, func() tea.Msg {
		ctx := context.Background()
		err := s.sessionService.Delete(ctx, selectedSession.ID)
		if err != nil {
			return utils.ReportError(err)
		}
		return SessionDeletedMsg{SessionID: selectedSession.ID}
	}
}

func (s *sessionDialogCmp) setupTextInput() {
	s.textInput = textinput.New()
	s.textInput.CharLimit = 100
	s.textInput.Width = 40
}

func (s *sessionDialogCmp) View() string {
	switch s.mode {
	case modeCreate:
		return s.renderInputMode("Create New Session", "Enter session title:")
	case modeEdit:
		return s.renderInputMode("Edit Session", "Edit session title:")
	case modeConfirmDelete:
		return s.renderDeleteConfirmation()
	default:
		return s.renderListMode()
	}
}

func (s *sessionDialogCmp) renderListMode() string {
	t := theme.CurrentTheme()
	baseStyle := styles.BaseStyle()

	if len(s.sessions) == 0 {
		return baseStyle.Padding(1, 2).
			Border(lipgloss.NormalBorder()).
			BorderBackground(t.Background()).
			BorderForeground(t.TextMuted()).
			Width(40).
			Render("No sessions available")
	}

	// Calculate max width needed for session titles
	maxWidth := 40 // Minimum width
	for _, sess := range s.sessions {
		if len(sess.Title) > maxWidth-4 { // Account for padding
			maxWidth = len(sess.Title) + 4
		}
	}

	maxWidth = max(30, min(maxWidth, s.width-15)) // Limit width to avoid overflow

	// Limit height to avoid taking up too much screen space
	maxVisibleSessions := min(10, len(s.sessions))

	// Build the session list
	sessionItems := make([]string, 0, maxVisibleSessions)
	startIdx := 0

	// If we have more sessions than can be displayed, adjust the start index
	if len(s.sessions) > maxVisibleSessions {
		// Center the selected item when possible
		halfVisible := maxVisibleSessions / 2
		if s.selectedIdx >= halfVisible && s.selectedIdx < len(s.sessions)-halfVisible {
			startIdx = s.selectedIdx - halfVisible
		} else if s.selectedIdx >= len(s.sessions)-halfVisible {
			startIdx = len(s.sessions) - maxVisibleSessions
		}
	}

	endIdx := min(startIdx+maxVisibleSessions, len(s.sessions))

	for i := startIdx; i < endIdx; i++ {
		sess := s.sessions[i]
		itemStyle := baseStyle.Width(maxWidth)

		if i == s.selectedIdx {
			itemStyle = itemStyle.
				Background(t.Primary()).
				Foreground(t.Background()).
				Bold(true)
		}

		sessionItems = append(sessionItems, itemStyle.Padding(0, 1).Render(sess.Title))
	}

	title := baseStyle.
		Foreground(t.Primary()).
		Bold(true).
		Width(maxWidth).
		Padding(0, 1).
		Render("Manage Sessions")

	// Add help text
	helpText := baseStyle.
		Foreground(t.TextMuted()).
		Width(maxWidth).
		Padding(0, 1).
		Render("n: new  e: edit  d: delete  enter: select  esc: close")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		baseStyle.Width(maxWidth).Render(""),
		baseStyle.Width(maxWidth).Render(lipgloss.JoinVertical(lipgloss.Left, sessionItems...)),
		baseStyle.Width(maxWidth).Render(""),
		helpText,
	)

	return baseStyle.Padding(1, 2).
		Border(lipgloss.NormalBorder()).
		BorderBackground(t.Background()).
		BorderForeground(t.TextMuted()).
		Width(lipgloss.Width(content) + 4).
		Render(content)
}

func (s *sessionDialogCmp) renderInputMode(title, prompt string) string {
	t := theme.CurrentTheme()
	baseStyle := styles.BaseStyle()
	
	titleText := baseStyle.
		Foreground(t.Primary()).
		Bold(true).
		Width(50).
		Padding(0, 1).
		Render(title)

	promptText := baseStyle.
		Width(50).
		Padding(0, 1).
		Render(prompt)

	inputView := s.textInput.View()
	
	helpText := baseStyle.
		Foreground(t.TextMuted()).
		Width(50).
		Padding(0, 1).
		Render("enter: save  esc: cancel")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		titleText,
		baseStyle.Width(50).Render(""),
		promptText,
		inputView,
		baseStyle.Width(50).Render(""),
		helpText,
	)

	return baseStyle.Padding(1, 2).
		Border(lipgloss.NormalBorder()).
		BorderBackground(t.Background()).
		BorderForeground(t.TextMuted()).
		Width(lipgloss.Width(content) + 4).
		Render(content)
}

func (s *sessionDialogCmp) renderDeleteConfirmation() string {
	t := theme.CurrentTheme()
	baseStyle := styles.BaseStyle()
	
	if len(s.sessions) == 0 {
		return s.renderListMode()
	}
	
	selectedSession := s.sessions[s.selectedIdx]
	
	titleText := baseStyle.
		Foreground(t.Primary()).
		Bold(true).
		Width(50).
		Padding(0, 1).
		Render("Delete Session")

	warningText := baseStyle.
		Foreground(lipgloss.Color("#ff6b6b")).
		Width(50).
		Padding(0, 1).
		Render(fmt.Sprintf("Delete session '%s'?", selectedSession.Title))

	confirmText := baseStyle.
		Width(50).
		Padding(0, 1).
		Render("This action cannot be undone.")
	
	helpText := baseStyle.
		Foreground(t.TextMuted()).
		Width(50).
		Padding(0, 1).
		Render("y: confirm  esc: cancel")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		titleText,
		baseStyle.Width(50).Render(""),
		warningText,
		confirmText,
		baseStyle.Width(50).Render(""),
		helpText,
	)

	return baseStyle.Padding(1, 2).
		Border(lipgloss.NormalBorder()).
		BorderBackground(t.Background()).
		BorderForeground(t.TextMuted()).
		Width(lipgloss.Width(content) + 4).
		Render(content)
}

func (s *sessionDialogCmp) BindingKeys() []key.Binding {
	return utils.KeyMapToSlice(sessionKeys)
}

func (s *sessionDialogCmp) SetSessionService(service session.Service) {
	s.sessionService = service
}

func (s *sessionDialogCmp) SetSessions(sessions []session.Session) {
	s.sessions = sessions

	// If we have a selected session ID, find its index
	if s.selectedSessionID != "" {
		for i, sess := range sessions {
			if sess.ID == s.selectedSessionID {
				s.selectedIdx = i
				return
			}
		}
	}

	// Default to first session if selected not found
	s.selectedIdx = 0
}

func (s *sessionDialogCmp) SetSelectedSession(sessionID string) {
	s.selectedSessionID = sessionID

	// Update the selected index if sessions are already loaded
	if len(s.sessions) > 0 {
		for i, sess := range s.sessions {
			if sess.ID == sessionID {
				s.selectedIdx = i
				return
			}
		}
	}
}

// NewSessionDialogCmp creates a new session switching dialog
func NewSessionDialogCmp() SessionDialog {
	s := &sessionDialogCmp{
		sessions:          []session.Session{},
		selectedIdx:       0,
		selectedSessionID: "",
		mode:              modeList,
	}
	s.setupTextInput()
	return s
}
