package page

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/yyovil/tandem/internal/app"
	"github.com/yyovil/tandem/internal/message"
	"github.com/yyovil/tandem/internal/session"
	"github.com/yyovil/tandem/internal/tools"
	"github.com/yyovil/tandem/internal/tui/bubbles/chat"
	"github.com/yyovil/tandem/internal/tui/layout"
	"github.com/yyovil/tandem/internal/utils"
)

var ChatPage PageID = "chat"

type chatPage struct {
	app      *app.App
	editor   layout.Container
	messages layout.Container
	layout   layout.SplitPaneLayout
	session  session.Session
}
type ChatKeyMap struct {
	NewSession key.Binding
	Cancel     key.Binding
}

var keyMap = ChatKeyMap{
	NewSession: key.NewBinding(
		key.WithKeys("ctrl+n"),
		key.WithHelp("ctrl+n", "new session"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
}

func (cp *chatPage) Init() tea.Cmd {
	cmds := []tea.Cmd{
		cp.layout.Init(),
	}
	return tea.Batch(cmds...)
}

func (cp *chatPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cmd := cp.layout.SetSize(msg.Width, msg.Height)
		cmds = append(cmds, cmd)
	case chat.SendMsg:
		cmd := cp.sendMessage(msg.Text, msg.Attachments)
		if cmd != nil {
			return cp, cmd
		}
	case chat.SessionSelectedMsg:
		if cp.session.ID == "" {
			cmd := cp.setSidebar()
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		cp.session = msg
	case tea.KeyMsg:
		switch {
		// Continue sending keys to layout->chat
		case key.Matches(msg, keyMap.NewSession):
			cp.session = session.Session{}
			return cp, tea.Batch(
				cp.clearSidebar(),
				utils.CmdHandler(chat.SessionClearedMsg{}),
			)
		case key.Matches(msg, keyMap.Cancel):
			if cp.session.ID != "" {
				// Cancel the current session's generation process
				// This allows users to interrupt long-running operations
				cp.app.Orchestrator.Cancel(cp.session.ID)
				return cp, nil
			}
		}
	}

	model, cmd := cp.layout.Update(msg)
	cmds = append(cmds, cmd)
	cp.layout = model.(layout.SplitPaneLayout)

	return cp, tea.Batch(cmds...)
}

func (cp *chatPage) setSidebar() tea.Cmd {
	sidebarContainer := layout.NewContainer(
		chat.NewSidebarCmp(cp.session),
	)
	return tea.Batch(cp.layout.SetRightPanel(sidebarContainer), sidebarContainer.Init())
}

func (cp *chatPage) clearSidebar() tea.Cmd {
	return cp.layout.ClearRightPanel()
}

// executeCommand executes a command using the terminal tool and returns the result
func (p *chatPage) executeCommand(command string) (string, error) {
	// Parse command into command and args
	parts := strings.Fields(strings.TrimSpace(command))
	if len(parts) == 0 {
		return "", fmt.Errorf("empty command")
	}

	cmd := parts[0]
	args := []string{}
	if len(parts) > 1 {
		args = parts[1:]
	}

	// Create terminal tool instance
	terminalTool := tools.NewDockerCli()

	// Prepare arguments for the terminal tool
	termArgs := tools.TerminalArgs{
		Command: cmd,
		Args:    args,
	}

	// Marshal the arguments to JSON
	argsJSON, err := json.Marshal(termArgs)
	if err != nil {
		return "", fmt.Errorf("failed to marshal terminal arguments: %w", err)
	}

	// Create tool call
	toolCall := tools.ToolCall{
		ID:    uuid.New().String(),
		Name:  "terminal",
		Input: string(argsJSON),
	}

	// Execute the command
	response, err := terminalTool.Run(context.Background(), toolCall)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %w", err)
	}

	if response.IsError {
		return "", fmt.Errorf("command execution failed: %s", response.Content)
	}

	return response.Content, nil
}

func (p *chatPage) sendMessage(text string, attachments []message.Attachment) tea.Cmd {
	var cmds []tea.Cmd
	if p.session.ID == "" {
		session, err := p.app.Sessions.Create(context.Background(), "New Session")
		if err != nil {
			return utils.ReportError(err)
		}

		p.session = session
		cmd := p.setSidebar()
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		cmds = append(cmds, utils.CmdHandler(chat.SessionSelectedMsg(session)))
	}

	// Check if this is a command execution request (starts with !)
	if strings.HasPrefix(text, "!") {
		// Extract command (everything after !)
		command := strings.TrimSpace(text[1:])
		if command == "" {
			return utils.ReportError(fmt.Errorf("empty command after !"))
		}

		// Execute the command using terminal tool
		result, err := p.executeCommand(command)
		if err != nil {
			return utils.ReportError(fmt.Errorf("command execution failed: %w", err))
		}

		// Create a user message showing the command
		_, err = p.app.Messages.Create(context.Background(), p.session.ID, message.CreateMessageParams{
			Role:  message.User,
			Parts: []message.ContentPart{message.TextContent{Text: text}},
		})
		if err != nil {
			return utils.ReportError(fmt.Errorf("failed to create user message: %w", err))
		}

		// Create an assistant message with the command result
		_, err = p.app.Messages.Create(context.Background(), p.session.ID, message.CreateMessageParams{
			Role:  message.Assistant,
			Parts: []message.ContentPart{message.TextContent{Text: fmt.Sprintf("```\n%s\n```", result)}},
		})
		if err != nil {
			return utils.ReportError(fmt.Errorf("failed to create result message: %w", err))
		}

		return tea.Batch(cmds...)
	}

	// Normal message processing - send to orchestrator
	_, err := p.app.Orchestrator.Run(context.Background(), p.session.ID, text, attachments...)
	if err != nil {
		return utils.ReportError(err)
	}
	return tea.Batch(cmds...)
}
func (p *chatPage) BindingKeys() []key.Binding {
	bindings := utils.KeyMapToSlice(keyMap)
	bindings = append(bindings, p.messages.BindingKeys()...)
	bindings = append(bindings, p.editor.BindingKeys()...)
	return bindings
}
func NewChatPage(app *app.App) tea.Model {

	messagesContainer := layout.NewContainer(
		chat.NewMessagesCmp(app),
	)
	editorContainer := layout.NewContainer(
		chat.NewEditorCmp(app),
		layout.WithBorder(true, false, false, false),
	)
	return &chatPage{
		app:      app,
		editor:   editorContainer,
		messages: messagesContainer,
		layout: layout.NewSplitPane(
			layout.WithLeftPanel(messagesContainer),
			layout.WithBottomPanel(editorContainer),
		),
	}
}

func (p *chatPage) View() string {
	layoutView := p.layout.View()
	return layoutView
}
