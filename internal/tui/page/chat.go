package page

import (
	"context"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yyovil/tandem/internal/app"
	"github.com/yyovil/tandem/internal/message"
	"github.com/yyovil/tandem/internal/session"
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
			// Only intercept ESC to cancel when the agent is actively working.
			if cp.session.ID != "" && cp.app.Orchestrator.IsSessionBusy(cp.session.ID) {
				cp.app.Orchestrator.Cancel(cp.session.ID)
				return cp, nil
			}
			// Otherwise, let ESC propagate to the editor (e.g., to exit EscapeShellMode).
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
