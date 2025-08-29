package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/yyovil/tandem/internal/app"
	"github.com/yyovil/tandem/internal/logging"
	"github.com/yyovil/tandem/internal/message"
	"github.com/yyovil/tandem/internal/session"
	"github.com/yyovil/tandem/internal/tools"
	"github.com/yyovil/tandem/internal/tui/bubbles/dialog"
	"github.com/yyovil/tandem/internal/tui/styles"
	"github.com/yyovil/tandem/internal/tui/theme"
	"github.com/yyovil/tandem/internal/utils"
)

type editorCmp struct {
	width           int
	height          int
	app             *app.App
	session         session.Session
	textarea        textarea.Model
	attachments     []message.Attachment
	EscapeShellMode bool
	deleteMode      bool
}

type EditorKeyMaps struct {
	EscapeShellCmd key.Binding
	Send           key.Binding
	OpenEditor     key.Binding
}

type bluredEditorKeyMaps struct {
	Send       key.Binding
	Focus      key.Binding
	OpenEditor key.Binding
}
type DeleteAttachmentKeyMaps struct {
	AttachmentDeleteMode key.Binding
	Escape               key.Binding
	DeleteAllAttachments key.Binding
}

var editorMaps = EditorKeyMaps{
	EscapeShellCmd: key.NewBinding(
		key.WithKeys("ctrl+enter"),
		key.WithHelp("ctrl+enter", "escape shell"),
	),
	Send: key.NewBinding(
		key.WithKeys("enter", "ctrl+s"),
		key.WithHelp("enter", "send message"),
	),
	OpenEditor: key.NewBinding(
		key.WithKeys("ctrl+e"),
		key.WithHelp("ctrl+e", "open editor"),
	),
}

var DeleteKeyMaps = DeleteAttachmentKeyMaps{
	AttachmentDeleteMode: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl+r+{i}", "delete attachment at index i"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel delete mode"),
	),
	DeleteAllAttachments: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("ctrl+r+r", "delete all attchments"),
	),
}

const (
	maxAttachments = 5
)

func (m *editorCmp) openEditor() tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nvim"
	}

	tmpfile, err := os.CreateTemp("", "msg_*.md")
	if err != nil {
		return utils.ReportError(err)
	}
	tmpfile.Close()
	c := exec.Command(editor, tmpfile.Name()) //nolint:gosec
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return utils.ReportError(err)
		}
		content, err := os.ReadFile(tmpfile.Name())
		if err != nil {
			return utils.ReportError(err)
		}
		if len(content) == 0 {
			return utils.ReportWarn("Message is empty")
		}
		os.Remove(tmpfile.Name())
		attachments := m.attachments
		m.attachments = nil
		return SendMsg{
			Text:        string(content),
			Attachments: attachments,
		}
	})
}

func (m *editorCmp) Init() tea.Cmd {
	return textarea.Blink
}

func (m *editorCmp) send() tea.Cmd {
	if m.app.Orchestrator.IsSessionBusy(m.session.ID) {
		return utils.ReportWarn("Agent is working, please wait...")
	}

	value := m.textarea.Value()
	m.textarea.Reset()
	attachments := m.attachments

	m.attachments = nil
	if value == "" {
		return nil
	}
	// If we're in EscapeShellMode, treat the input as a shell command
	if m.EscapeShellMode {
		// Expect prefix "! " to indicate shell escape
		cmdline := strings.TrimSpace(strings.TrimPrefix(value, "! "))
		if cmdline == "" {
			return nil
		}
		return m.EscapeShell(cmdline)
	}
	return tea.Batch(
		utils.CmdHandler(SendMsg{
			Text:        value,
			Attachments: attachments,
		}),
	)
}

func (m *editorCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case SessionSelectedMsg:
		if msg.ID != m.session.ID {
			m.session = msg
		}
		return m, nil
	case dialog.AttachmentAddedMsg:
		if len(m.attachments) >= maxAttachments {
			logging.ErrorPersist(fmt.Sprintf("cannot add more than %d images", maxAttachments))
			return m, cmd
		}
		m.attachments = append(m.attachments, msg.Attachment)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DeleteKeyMaps.AttachmentDeleteMode):
			m.deleteMode = true
			return m, nil
		case key.Matches(msg, DeleteKeyMaps.DeleteAllAttachments) && m.deleteMode:
			m.deleteMode = false
			m.attachments = nil
			return m, nil
		case m.deleteMode && len(msg.Runes) > 0 && unicode.IsDigit(msg.Runes[0]):
			num := int(msg.Runes[0] - '0')
			m.deleteMode = false
			if num < 10 && len(m.attachments) > num {
				if num == 0 {
					m.attachments = m.attachments[num+1:]
				} else {
					m.attachments = slices.Delete(m.attachments, num, num+1)
				}
				return m, nil
			}
		case key.Matches(msg, messageKeys.PageUp) || key.Matches(msg, messageKeys.PageDown) ||
			key.Matches(msg, messageKeys.HalfPageUp) || key.Matches(msg, messageKeys.HalfPageDown):
			return m, nil
		case key.Matches(msg, editorMaps.OpenEditor):
			if m.app.Orchestrator.IsSessionBusy(m.session.ID) {
				return m, utils.ReportWarn("Agent is working, please wait...")
			}
			return m, m.openEditor()
		case key.Matches(msg, DeleteKeyMaps.Escape):
			m.deleteMode = false
			// Exit EscapeShellMode on ESC and remove leading "! " if present
			if m.EscapeShellMode {
				m.EscapeShellMode = false
				val := m.textarea.Value()
				if strings.HasPrefix(val, "! ") {
					m.textarea.SetValue(strings.TrimPrefix(val, "! "))
				}
			}
			return m, nil
		case m.textarea.Focused() && key.Matches(msg, editorMaps.Send):
			// Handle Enter key
			value := m.textarea.Value()
			if len(value) > 0 && value[len(value)-1] == '\\' {
				// If the last character is a backslash, remove it and add a newline
				m.textarea.SetValue(value[:len(value)-1] + "\n")
				return m, nil
			} else {
				// Otherwise, send the message
				return m, m.send()
			}
		}

	}
	// First update textarea to capture latest input, then update modes/prompts
	m.textarea, cmd = m.textarea.Update(msg)
	// Enter EscapeShellMode when input starts with "! ". Trim the trigger from the textarea value.
	val := m.textarea.Value()
	if !m.EscapeShellMode && strings.HasPrefix(val, "! ") {
		m.EscapeShellMode = true
		m.textarea.SetValue(strings.TrimPrefix(val, "! "))
	}
	return m, cmd
}

func (m *editorCmp) View() string {
	t := theme.CurrentTheme()

	// Style the prompt with theme colors
	style := lipgloss.NewStyle().
		Bold(true).
		PaddingLeft(2).
		Height(m.textarea.Height()).
		Background(t.Background()).
		Foreground(t.Primary())

	prompt := ">"
	if m.EscapeShellMode {
		prompt = "!"
		style = style.Foreground(t.Secondary())
	}

	if len(m.attachments) == 0 {
		return lipgloss.JoinHorizontal(lipgloss.Top, style.Render(prompt), m.textarea.View())
	}

	m.textarea.SetHeight(m.height - 1)
	return lipgloss.JoinVertical(lipgloss.Top,
		m.attachmentsContent(),
		lipgloss.JoinHorizontal(lipgloss.Top, style.Render(prompt),
			m.textarea.View()),
	)
}

func (m *editorCmp) SetSize(width, height int) tea.Cmd {
	m.width = width
	m.height = height
	m.textarea.SetWidth(width - 3) // account for the prompt + ">"
	m.textarea.SetHeight(height)
	return nil
}

func (m *editorCmp) GetSize() (int, int) {
	return m.textarea.Width(), m.textarea.Height()
}

func (m *editorCmp) attachmentsContent() string {
	var styledAttachments []string
	t := theme.CurrentTheme()
	attachmentStyles := styles.BaseStyle().
		MarginLeft(1).
		MarginBackground(t.Background()).
		Background(t.BackgroundSecondary()).
		Foreground(t.Text())

	for i, attachment := range m.attachments {
		var filename string
		if len(attachment.FileName) > 10 {
			filename = fmt.Sprintf(" %s %s...", styles.DocumentIcon, attachment.FileName[0:7])
		} else {
			filename = fmt.Sprintf(" %s %s", styles.DocumentIcon, attachment.FileName)
		}
		if m.deleteMode {
			filename = fmt.Sprintf("%d%s", i, filename)
		}
		styledAttachments = append(styledAttachments, attachmentStyles.Render(filename))
	}
	content := styles.BaseStyle().Width(m.width).Render(lipgloss.JoinHorizontal(lipgloss.Left, styledAttachments...))

	return content
}

func (m *editorCmp) BindingKeys() []key.Binding {
	bindings := []key.Binding{}
	bindings = append(bindings, utils.KeyMapToSlice(editorMaps)...)
	bindings = append(bindings, utils.KeyMapToSlice(DeleteKeyMaps)...)
	return bindings
}

func CreateTextArea(existing *textarea.Model) textarea.Model {
	t := theme.CurrentTheme()
	bgColor := t.Background()
	textColor := t.Text()
	textMutedColor := t.TextMuted()

	ta := textarea.New()
	ta.BlurredStyle.Base = styles.BaseStyle().Background(bgColor).Foreground(textColor)
	ta.BlurredStyle.CursorLine = styles.BaseStyle().Background(bgColor)
	ta.BlurredStyle.Placeholder = styles.BaseStyle().Background(bgColor).Foreground(textMutedColor)
	ta.BlurredStyle.Text = styles.BaseStyle().Background(bgColor).Foreground(textColor)
	ta.FocusedStyle.Base = styles.BaseStyle().Background(bgColor).Foreground(textColor)
	ta.FocusedStyle.CursorLine = styles.BaseStyle().Background(bgColor)
	ta.FocusedStyle.Placeholder = styles.BaseStyle().Background(bgColor).Foreground(textMutedColor)
	ta.FocusedStyle.Text = styles.BaseStyle().Background(bgColor).Foreground(textColor)

	ta.Prompt = " "
	ta.ShowLineNumbers = false
	ta.CharLimit = -1

	if existing != nil {
		ta.SetValue(existing.Value())
		ta.SetWidth(existing.Width())
		ta.SetHeight(existing.Height())
	}

	ta.Focus()
	return ta
}

func NewEditorCmp(app *app.App) tea.Model {
	ta := CreateTextArea(nil)
	return &editorCmp{
		app:      app,
		textarea: ta,
	}
}

// EscapeShell executes a shell command inside the Kali Docker container and streams results into the chat as tool output.
func (m *editorCmp) EscapeShell(cmdline string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		// Ensure session exists
		if m.session.ID == "" {
			sess, err := m.app.Sessions.Create(ctx, "Shell Session")
			if err != nil {
				return utils.InfoMsg{Type: utils.InfoTypeError, Msg: err.Error()}
			}
			m.session = sess
			// Inform rest of the app about the new session
			// Returning a SessionSelectedMsg won't prevent pubsub updates
			// but we can still publish it by returning it after work below if needed.
		}

		// Create a user message reflecting the command
		userText := "! " + cmdline
		_, _ = m.app.Messages.Create(ctx, m.session.ID, message.CreateMessageParams{
			Role:  message.User,
			Parts: []message.ContentPart{message.TextContent{Text: userText}},
		})

		// Prepare tool call for terminal using the user's command and args as-is
		fields := strings.Fields(cmdline)
		if len(fields) == 0 {
			return SessionSelectedMsg(m.session)
		}
		termArgs := tools.TerminalArgs{Command: fields[0]}
		if len(fields) > 1 {
			termArgs.Args = fields[1:]
		}
		payload, _ := json.Marshal(termArgs)
		callID := uuid.New().String()

		// Create assistant message with the tool call (marked finished so UI shows params + response)
		assistantParts := []message.ContentPart{
			message.ToolCall{ID: callID, Name: tools.TerminalToolName, Input: string(payload), Type: "", Finished: true},
		}
		_, _ = m.app.Messages.Create(ctx, m.session.ID, message.CreateMessageParams{
			Role:  message.Assistant,
			Parts: assistantParts,
		})

		// Run the tool (reuse ExecuteCmd for execution) without shell wrapping
		tool := tools.NewDockerCli().(*tools.Terminal)
		output, err := tool.ExecuteCmd(ctx, fields)
		resp := tools.ToolResponse{Type: tools.ToolResponseTypeText, Content: output}
		if err != nil {
			resp = tools.NewTextErrorResponse("terminal execution failed: " + err.Error())
		}

		// Create tool result message
		_, _ = m.app.Messages.Create(ctx, m.session.ID, message.CreateMessageParams{
			Role: message.Tool,
			Parts: []message.ContentPart{message.ToolResult{
				ToolCallID: callID,
				Content:    resp.Content,
				Metadata:   resp.Metadata,
				IsError:    resp.IsError,
			}},
		})

		// If we had to create a new session locally, notify the rest of the app
		return SessionSelectedMsg(m.session)
	}
}
