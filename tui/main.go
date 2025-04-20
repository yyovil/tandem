package main

import (
	"fmt"
	// "log"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"os"
)

type State int

const (
	// when user be writing the prompt. from this state, app can only go to requesting state.
	prompting State = iota
	// when a completion stream is requested. from this state, app can only go to interrupted state.
	requesting
	// when the response is streaming. from this state, app can only go to prompting state or interrupted.
	streaming
	// when streaming failed or was aborted
	interrupted
)

type streamChunkMsg struct{ Text string }
type streamDoneMsg struct{}

// this model reprs your entire state of the cli app.
type Model struct {
	state    State
	viewport viewport.Model
	textarea textarea.Model
	spinner  spinner.Model
	messages []string
	message,
	sessionId,
	userId,
	streamingResponse,
	/*
		TODO: user should be able to select the model by himself because you don't know how they be feeling some type of way. create a ENUM for models.
	*/
	model string
}

func (m *Model) Init() tea.Cmd {

	return textarea.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 10
		m.viewport.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).BorderStyle(lipgloss.NormalBorder())

		m.textarea.SetWidth(msg.Width)
		m.textarea.SetHeight(msg.Height / 10)

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		// continue ticking if still requesting.
		if m.state == requesting {
			return m, cmd
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			m.state = interrupted
			return m, tea.Quit
		case tea.KeyTab:
			m.state = requesting
			m.message = m.textarea.Value()
			m.messages = append(m.messages, "You: "+m.textarea.Value())
			return m, tea.Sequence(m.spinner.Tick, m.GetCompletionStreamCmd())
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}
	case streamChunkMsg:

		if m.state == requesting {
			m.state = streaming
		}

		// 1) accumulate raw markdown
		m.streamingResponse += msg.Text

		// 2) render it via Glamour
		rendered, err := glamour.Render(m.streamingResponse, "dark")
		if err != nil {
			m.viewport.SetContent(m.streamingResponse)
		} else {
			m.viewport.SetContent(rendered)
		}

		// 3) keep pulling more chunks
		return m, nil

	case streamDoneMsg:
		m.state = prompting
		m.textarea.Reset()
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// update the View to render the streaming response.
func (m *Model) View() string {
	tuiLayout := "Chat with Sage.\n\n%s\n\n%s\n\n%s %s"
	return fmt.Sprintf(
		tuiLayout,
		m.viewport.View(),
		m.textarea.View(),
		m.spinner.View(),
		"ctrl+c: quit | tab: to send",
	) + "\n\n"
}

func initialModel() tea.Model {
	s := spinner.New()
	s.Spinner = spinner.Meter

	vp := viewport.New(1, 1)
	vp.Style = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))
	ta := textarea.New()
	ta.FocusedStyle = textarea.Style{
		Base: lipgloss.NewStyle().Border(lipgloss.NormalBorder()),
	}
	ta.ShowLineNumbers = false

	ta.Placeholder = "Enter your prompt..."
	ta.Focus()

	return &Model{
		state:             prompting,
		spinner:           s,
		streamingResponse: "",
		viewport:          vp,
		textarea:          ta,
		messages:          []string{},
		message:           "",
		model:             ModelGemini20FlashLite,
		userId:            "slimeMaster",
		sessionId:         "slimeMasterSession1",
	}
}

func main() {
	_, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		os.Exit(2)
	}
	initialModel()

	tuiLoop := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := tuiLoop.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	// endpoint := os.Getenv("ENDPOINT")
	// if endpoint != "" {
	// 	log.Println("oh fuck! ENDPOINT var is either not (defined or exported).")
	// }
	// CreateSSEClient(endpoint)
}

/*
TODO:
	>  use errors pkgs to provide better error handling instead of trying to log to stdout as that is being occupied by the tui.
*/
