package bubbles

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yyovil/tandem/internal/config"
	"github.com/yyovil/tandem/internal/models"
	"github.com/yyovil/tandem/internal/pubsub"
	"github.com/yyovil/tandem/internal/session"
	"github.com/yyovil/tandem/internal/tui/bubbles/chat"
	"github.com/yyovil/tandem/internal/tui/styles"
	"github.com/yyovil/tandem/internal/tui/theme"
	"github.com/yyovil/tandem/internal/utils"
)

type StatusCmp interface {
	tea.Model
}

type statusCmp struct {
	info       utils.InfoMsg
	width      int
	messageTTL time.Duration
	session    session.Session
}

// clearMessageCmd is a command that clears status messages after a timeout
func (m statusCmp) clearMessageCmd(ttl time.Duration) tea.Cmd {
	return tea.Tick(ttl, func(time.Time) tea.Msg {
		return utils.ClearStatusMsg{}
	})
}

func (m statusCmp) Init() tea.Cmd {
	return nil
}

func (m statusCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	case chat.SessionSelectedMsg:
		m.session = msg
	case chat.SessionClearedMsg:
		m.session = session.Session{}
	case pubsub.Event[session.Session]:
		if msg.Type == pubsub.UpdatedEvent {
			if m.session.ID == msg.Payload.ID {
				m.session = msg.Payload
			}
		}
	case utils.InfoMsg:
		m.info = msg
		ttl := msg.TTL
		if ttl == 0 {
			ttl = m.messageTTL
		}
		return m, m.clearMessageCmd(ttl)
	case utils.ClearStatusMsg:
		m.info = utils.InfoMsg{}
	}
	return m, nil
}

var helpWidget = ""

// getHelpWidget returns the help widget with current theme colors
func getHelpWidget() string {
	t := theme.CurrentTheme()
	helpText := "ctrl+? help"

	return styles.Padded().
		Background(t.TextMuted()).
		Foreground(t.BackgroundDarker()).
		Bold(true).
		Render(helpText)
}

func formatTokensAndCost(tokens, contextWindow int64, cost float64) string {
	// Format tokens in human-readable format (e.g., 110K, 1.2M)
	var formattedTokens string
	switch {
	case tokens >= 1_000_000:
		formattedTokens = fmt.Sprintf("%.1fM", float64(tokens)/1_000_000)
	case tokens >= 1_000:
		formattedTokens = fmt.Sprintf("%.1fK", float64(tokens)/1_000)
	default:
		formattedTokens = fmt.Sprintf("%d", tokens)
	}

	// Remove .0 suffix if present
	if strings.HasSuffix(formattedTokens, ".0K") {
		formattedTokens = strings.Replace(formattedTokens, ".0K", "K", 1)
	}
	if strings.HasSuffix(formattedTokens, ".0M") {
		formattedTokens = strings.Replace(formattedTokens, ".0M", "M", 1)
	}

	// Format cost with $ symbol and 2 decimal places
	formattedCost := fmt.Sprintf("$%.2f", cost)

	percentage := (float64(tokens) / float64(contextWindow)) * 100
	if percentage > 80 {
		// add the warning icon and percentage
		formattedTokens = fmt.Sprintf("%s(%d%%)", styles.WarningIcon, int(percentage))
	}

	return fmt.Sprintf("Context: %s, Cost: %s", formattedTokens, formattedCost)
}

func (m statusCmp) View() string {
	t := theme.CurrentTheme()
	modelID := config.Get().Agents[config.Orchestrator].Model
	model := models.SupportedModels[modelID]

	// Initialize the help widget
	status := getHelpWidget()

	tokenInfoWidth := 0
	if m.session.ID != "" {
		totalTokens := m.session.PromptTokens + m.session.CompletionTokens
		tokens := formatTokensAndCost(totalTokens, model.ContextWindow, m.session.Cost)
		tokensStyle := styles.Padded().
			Background(t.Text()).
			Foreground(t.BackgroundSecondary())
		percentage := (float64(totalTokens) / float64(model.ContextWindow)) * 100
		if percentage > 80 {
			tokensStyle = tokensStyle.Background(t.Warning())
		}
		tokenInfoWidth = lipgloss.Width(tokens) + 2
		status += tokensStyle.Render(tokens)
	}

	availableWidth := max(0, m.width-lipgloss.Width(helpWidget)-lipgloss.Width(m.model())-tokenInfoWidth)

	if m.info.Msg != "" {
		infoStyle := styles.Padded().
			Foreground(t.Background()).
			Width(availableWidth)

		switch m.info.Type {
		case utils.InfoTypeInfo:
			infoStyle = infoStyle.Background(t.Info())
		case utils.InfoTypeWarn:
			infoStyle = infoStyle.Background(t.Warning())
		case utils.InfoTypeError:
			infoStyle = infoStyle.Background(t.Error())
		}

		infoWidth := availableWidth - 10
		// Truncate message if it's longer than available width
		msg := m.info.Msg
		if len(msg) > infoWidth && infoWidth > 0 {
			msg = msg[:infoWidth] + "..."
		}
		status += infoStyle.Render(msg)
	} else {
		status += styles.Padded().
			Foreground(t.Text()).
			Background(t.BackgroundSecondary()).
			Width(availableWidth).
			Render("")
	}

	status += m.model()
	return status
}


func (m statusCmp) model() string {
	t := theme.CurrentTheme()

	cfg := config.Get()

	orchestrator, ok := cfg.Agents[config.Orchestrator]
	if !ok {
		return "Unknown"
	}
	model := models.SupportedModels[orchestrator.Model]

	return styles.Padded().
		Background(t.Secondary()).
		Foreground(t.Background()).
		Render(model.Name)
}

func NewStatusCmp() StatusCmp {
	helpWidget = getHelpWidget()

	return &statusCmp{
		messageTTL: 10 * time.Second,
	}
}
