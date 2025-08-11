package page

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yaydraco/tandem/internal/app"
	"github.com/yaydraco/tandem/internal/pubsub"
	"github.com/yaydraco/tandem/internal/subagent"
	"github.com/yaydraco/tandem/internal/tui/theme"
)

var ActivityPage PageID = "activity"

type activityPageModel struct {
	app    *app.App
	table  table.Model
	width  int
	height int
	keys   activityKeyMap
}

type activityKeyMap struct {
	Abort key.Binding
	Refresh key.Binding
}

var activityKeys = activityKeyMap{
	Abort: key.NewBinding(
		key.WithKeys("ctrl+a"),
		key.WithHelp("ctrl+a", "abort selected task"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
}

type AbortActivityMsg struct {
	ActivityID string
}

func NewActivityPage(app *app.App) tea.Model {
	t := theme.CurrentTheme()
	
	columns := []table.Column{
		{Title: "Agent", Width: 15},
		{Title: "Task", Width: 25},
		{Title: "Status", Width: 20},
		{Title: "Progress", Width: 8},
		{Title: "ETA", Width: 8},
		{Title: "Started", Width: 8},
		{Title: "Duration", Width: 8},
	}
	
	tbl := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)
	
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(t.BorderFocused()).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(t.Text()).
		Background(t.BorderFocused()).
		Bold(false)
	tbl.SetStyles(s)
	
	return &activityPageModel{
		app:   app,
		table: tbl,
		keys:  activityKeys,
	}
}

func (m *activityPageModel) Init() tea.Cmd {
	return m.refreshActivities()
}

func (m *activityPageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		// Update table size
		m.table.SetWidth(msg.Width - 4)
		m.table.SetHeight(msg.Height - 8)
		
		// Adjust column widths based on available space
		totalWidth := msg.Width - 8 // Leave some margin
		cols := []table.Column{
			{Title: "Agent", Width: totalWidth * 15 / 100},
			{Title: "Task", Width: totalWidth * 25 / 100},
			{Title: "Status", Width: totalWidth * 20 / 100},
			{Title: "Progress", Width: totalWidth * 8 / 100},
			{Title: "ETA", Width: totalWidth * 8 / 100},
			{Title: "Started", Width: totalWidth * 12 / 100},
			{Title: "Duration", Width: totalWidth * 12 / 100},
		}
		m.table.SetColumns(cols)
		
		return m, nil
		
	case pubsub.Event[subagent.ActivityEvent]:
		// Refresh the table when activities update
		return m, m.refreshActivities()
		
	case AbortActivityMsg:
		if m.app.SubAgents != nil {
			m.app.SubAgents.AbortActivity(context.Background(), msg.ActivityID)
		}
		return m, nil
		
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Abort):
			if m.table.Cursor() < len(m.table.Rows()) {
				selected := m.table.SelectedRow()
				if len(selected) > 0 {
					// The activity ID is stored in a hidden column or we can extract from the row
					// For now, we'll need to get it from the activities list
					activities := m.getActivities()
					if m.table.Cursor() < len(activities) {
						activity := activities[m.table.Cursor()]
						if activity.CanAbort {
							return m, func() tea.Msg {
								return AbortActivityMsg{ActivityID: activity.ID}
							}
						}
					}
				}
			}
			return m, nil
			
		case key.Matches(msg, m.keys.Refresh):
			return m, m.refreshActivities()
		}
	}
	
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)
	
	return m, tea.Batch(cmds...)
}

func (m *activityPageModel) View() string {
	t := theme.CurrentTheme()
	
	titleStyle := lipgloss.NewStyle().
		Foreground(t.Primary()).
		Bold(true).
		Padding(0, 1)
	title := titleStyle.Render("🤖 SubAgent Activities")
	
	var content string
	activities := m.getActivities()
	
	if len(activities) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(t.TextMuted()).
			Align(lipgloss.Center).
			Height(m.height - 6)
		content = emptyStyle.Render("No active subagent tasks")
	} else {
		content = m.table.View()
	}
	
	help := lipgloss.NewStyle().
		Foreground(t.TextMuted()).
		Render("ctrl+a: abort task • r: refresh • esc: back")
	
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		content,
		"",
		help,
	)
}

func (m *activityPageModel) refreshActivities() tea.Cmd {
	return func() tea.Msg {
		activities := m.getActivities()
		rows := make([]table.Row, len(activities))
		
		for i, activity := range activities {
			// Format duration
			duration := time.Since(activity.StartedAt).Truncate(time.Second)
			
			// Truncate task description if too long
			task := activity.Task
			if len(task) > 20 {
				task = task[:17] + "..."
			}
			
			// Use the enhanced status text
			statusText := activity.StatusText
			if len(statusText) > 18 {
				statusText = statusText[:15] + "..."
			}
			
			// Format status with color indicators
			status := statusText
			switch activity.Status {
			case subagent.StatusStarting:
				status = "🔄 " + statusText
			case subagent.StatusRunning:
				status = "⚡ " + statusText
			case subagent.StatusCompleted:
				status = "✅ " + statusText
			case subagent.StatusError:
				status = "❌ " + statusText
			case subagent.StatusAborted:
				status = "🛑 " + statusText
			}
			
			// Truncate status if too long
			if len(status) > 18 {
				status = status[:15] + "..."
			}
			
			estimatedTime := activity.EstimatedTime
			if estimatedTime == "" {
				estimatedTime = "-"
			}
			
			rows[i] = table.Row{
				string(activity.AgentName),
				task,
				status,
				activity.Progress,
				estimatedTime,
				activity.StartedAt.Format("15:04:05"),
				duration.String(),
			}
		}
		
		m.table.SetRows(rows)
		return nil
	}
}

func (m *activityPageModel) getActivities() []subagent.Activity {
	if m.app.SubAgents == nil {
		return []subagent.Activity{}
	}
	
	return m.app.SubAgents.GetActiveActivities(context.Background())
}

func (m *activityPageModel) BindingKeys() []key.Binding {
	return []key.Binding{
		m.keys.Abort,
		m.keys.Refresh,
	}
}

func (m *activityPageModel) SetSize(width, height int) tea.Cmd {
	m.width = width
	m.height = height
	return nil
}