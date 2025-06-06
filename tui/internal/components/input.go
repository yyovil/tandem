package components

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	vp "github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yyovil/tui/internal/components/messages"
	"github.com/yyovil/tui/internal/utils"
)

type Status string

const (
	Streaming  Status = "Streaming"
	Requesting Status = "Requesting"
	Idle       Status = "Idle"
)

type Input struct {
	status        Status
	stream        chan tea.Msg
	width, height int
	userPrompt    string
	textarea      textarea.Model
	FilePicker    FilePicker
	// TODO: out this and put in a dedicated layout file.
	leftpane, rightpane vp.Model

	leftPaneMessages []tea.Msg //TODO: this out, use a better type
}

type InputKeyMap struct {
	ShowFilePicker,
	Send,
	Quit,
	PageDown,
	PageUp,
	HalfPageUp,
	HalfPageDown key.Binding
}

var inputKeyMap = InputKeyMap{
	ShowFilePicker: key.NewBinding(
		key.WithKeys("ctrl+o"),
		key.WithHelp("ctrl+o", "attach file"),
	),
	Send: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "send message"),
	),
	Quit: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "quit"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("f/pgdn", "page down"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("b/pgupf", "page up"),
	),
	HalfPageUp: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("up", "½ page up"),
	),
	HalfPageDown: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("down", "½ page down"),
	),
}

func (i *Input) Init() tea.Cmd {
	i.textarea.Placeholder = "Assign tasks to AI Agents here..."
	i.textarea.Focus()
	i.textarea.ShowLineNumbers = false
	return textarea.Blink
}

func (i *Input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, inputKeyMap.PageUp),
			key.Matches(msg, inputKeyMap.HalfPageUp),
			key.Matches(msg, inputKeyMap.PageDown),
			key.Matches(msg, inputKeyMap.HalfPageDown):
			break

		case key.Matches(msg, inputKeyMap.ShowFilePicker):
			cmd = i.FilePicker.Init()
			cmds = append(cmds, cmd)
			i.FilePicker.showFilePicker = true
		case key.Matches(msg, inputKeyMap.Send):
			if !i.FilePicker.showFilePicker {

				if i.textarea.Value() == "" || i.status != Idle {
					return i, nil
				}

				i.userPrompt = i.textarea.Value()

				cmds = append(cmds, i.sendRunRequestCmd(), messages.AddUserMessageCmd(i.userPrompt, i.FilePicker.selectedFiles))

				i.textarea.Reset()
				i.FilePicker.viewport.GotoTop()
				i.FilePicker.filepicker.FileSelected = ""
				i.FilePicker.selectedFiles = nil

				return i, tea.Batch(cmds...)
			}

			_, cmd = i.FilePicker.Update(msg)
			cmds = append(cmds, cmd)
			return i, tea.Batch(cmds...)

		case key.Matches(msg, inputKeyMap.Quit):
			if !i.FilePicker.showFilePicker {
				return i, tea.Quit
			}

		default:
			cmd = i.textarea.Focus()
			cmds = append(cmds, cmd)
			i.textarea, cmd = i.textarea.Update(msg)
			cmds = append(cmds, cmd)
			return i, tea.Batch(cmds...)
		}
	case tea.WindowSizeMsg:
		i.textarea.SetWidth(msg.Width - 2)
		i.textarea.MaxWidth = msg.Width - 2
		i.textarea.SetHeight(4)
		i.textarea.MaxHeight = 5

		i.height = msg.Height
		i.width = msg.Width

		leftPaneWidth := ((i.width * 70) / 100) - 1
		rightPaneWidth := (i.width * 30) / 100
		paneHeight := i.height - 10

		i.leftpane.Width = leftPaneWidth
		i.leftpane.Height = paneHeight

		i.rightpane.Width = rightPaneWidth
		i.rightpane.Height = paneHeight

	case messages.UserMessageAddedMsg:
		msg.UserMessage.Width = i.leftpane.Width
		i.leftPaneMessages = append(i.leftPaneMessages, msg.UserMessage)

	case messages.AgentMessageAddedMsg:
		// blocking call to receive the first chunk of the stream
		agentMessage := messages.AgentMessage{
			StreamChan: msg.StreamChan,
			Width:      i.leftpane.Width,
		}
		firstChunk, _ := <-msg.StreamChan
		if v, ok := firstChunk.(messages.ConcatenateChunkMsg); ok {
			agentMessage.Content = string(v)
			i.leftPaneMessages = append(i.leftPaneMessages, agentMessage)
			i.leftpane.GotoBottom()

			return i, messages.ListenOnStreamChanCmd(msg.StreamChan)
		}

	case messages.ConcatenateChunkMsg:
		// Update the last message in place instead of appending a new one
		if len(i.leftPaneMessages) > 0 {
			lastMsgIndex := len(i.leftPaneMessages) - 1
			lastMsg := i.leftPaneMessages[lastMsgIndex]

			if agentMsg, ok := lastMsg.(messages.AgentMessage); ok {
				agentMsg.Content += string(msg)
				i.leftPaneMessages[lastMsgIndex] = agentMsg
				i.leftpane.GotoBottom()

				return i, messages.ListenOnStreamChanCmd(agentMsg.StreamChan)
			}
		}

	case messages.EndStream:
		i.status = Idle
		return i, nil
	}

	i.leftpane, cmd = i.leftpane.Update(msg)
	cmds = append(cmds, cmd)
	_, cmd = i.FilePicker.Update(msg)
	cmds = append(cmds, cmd)
	i.textarea, cmd = i.textarea.Update(msg)
	cmds = append(cmds, cmd)
	i.rightpane, cmd = i.rightpane.Update(msg)
	cmds = append(cmds, cmd)

	return i, tea.Batch(cmds...)
}

func (i *Input) View() string {

	if i.FilePicker.showFilePicker {
		return i.FilePicker.View()
	}

	inputStyle := lipgloss.
		NewStyle().
		Width(i.width-2).
		MaxWidth(i.width).
		Height(4).
		MaxHeight(6).
		Border(lipgloss.NormalBorder(), true)

	leftPaneStyle := lipgloss.NewStyle().
		Width(((i.width*70)/100)-1).
		MaxWidth((i.width*70)/100).
		Height(i.height-inputStyle.GetHeight()-10).
		MaxHeight(i.height-inputStyle.GetHeight()).
		// Background(lipgloss.Color("#d67ab1")).
		// Border(lipgloss.InnerHalfBlockBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color("#d67ab1")).
		Padding(1, 0)

	rightPaneStyle := lipgloss.NewStyle().
		Width(((i.width * 30) / 100)).
		//FIX: this is weird hack to get the width of rightpanel correct.
		MaxWidth((i.width*30)/100+1).
		Height(i.height-inputStyle.GetHeight()-10).
		MaxHeight(i.height-inputStyle.GetHeight()).
		// Border(lipgloss.InnerHalfBlockBorder(), false, false, false, true).
		// BorderLeftForeground(lipgloss.Color("#e2a3c7")).
		// Background(lipgloss.Color("#e2a3c7")).
		Padding(1, 0)

	content := ""
	for _, msg := range i.leftPaneMessages {
		switch m := msg.(type) {
		case messages.AgentMessage:
			content += m.View() + "\n\n"
		case messages.UserMessage:
			content += m.View() + "\n\n"
		}
	}

	i.leftpane.SetContent(content)

	panes := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPaneStyle.Render(i.leftpane.View()),
		rightPaneStyle.Render(i.rightpane.View()),
	)

	return lipgloss.JoinVertical(lipgloss.Top, panes, inputStyle.Render(i.textarea.View())+"\n"+i.footerView())
}

func (i Input) footerView() string {
	var s strings.Builder
	footerStyle := lipgloss.
		NewStyle().
		Height(1).
		Width(i.width-12).
		MaxHeight(2).
		Border(lipgloss.InnerHalfBlockBorder(), false).
		BorderLeft(true).
		PaddingLeft(1).
		Background(lipgloss.Color("#343a40")).
		MarginBottom(1)

	fpSelectedStyle := i.FilePicker.filepicker.Styles.Selected
	selectedFiles := i.FilePicker.selectedFiles

	if len(selectedFiles) == 0 {
		s.WriteString("No attachments")
	} else if len(selectedFiles) == 1 {
		footerStyle = footerStyle.BorderForeground(lipgloss.Color("212"))
		s.WriteString("Selected file: " + fpSelectedStyle.Render(selectedFiles[0]))
	} else {
		footerStyle = footerStyle.BorderForeground(lipgloss.Color("212"))
		s.WriteString("Total attachments: " + fpSelectedStyle.Render(fmt.Sprintf("%d", len(selectedFiles))))
	}

	statusStyle := lipgloss.
		NewStyle().
		Width(10).
		Height(1).
		Border(lipgloss.InnerHalfBlockBorder(), false).
		BorderRight(true).
		Background(lipgloss.Color("#fb5607")).
		AlignHorizontal(lipgloss.Center)

	return lipgloss.JoinHorizontal(lipgloss.Left, footerStyle.Render(s.String()), statusStyle.Render(string(i.status)))
}

func NewInput() Input {
	return Input{
		status:           Idle,
		leftPaneMessages: []tea.Msg{},
		width:            0,
		height:           0,
		userPrompt:       "",
		textarea:         textarea.New(),
		FilePicker:       NewFilePicker(),
		leftpane:         vp.New(0, 0),
		rightpane:        vp.New(0, 0),
	}
}

func (i *Input) sendRunRequestCmd() tea.Cmd {
	i.status = Requesting

	return func() tea.Msg {
		attachment, err := i.FilePicker.GetSelectedFile()
		if err != nil {
			i.status = Idle
			log.Println("error getting the attachment:", err.Error())
			return func() tea.Msg {
				return nil
			}
		}

		req, err := utils.GetPostRequest(i.userPrompt, attachment)
		if err != nil {
			i.status = Idle
			log.Println("error creating request:", err.Error())
			return func() tea.Msg {
				return nil
			}
		}

		stream := make(chan tea.Msg)

		go func() {
			defer close(stream)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				i.status = Idle
				log.Println("error sending request:", err.Error())
				// TODO: show a user feedback for this error
				log.Println("input status is: ", i.status)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				i.status = Idle
				log.Println("error: received non-200 response:", resp.Status)
				// TODO: show a user feedback for this error
				return
			}

			scanner := bufio.NewScanner(resp.Body)
			scanner.Split(bufio.ScanRunes)
			for scanner.Scan() {
				chunk := scanner.Text()
				if chunk != "" {
					stream <- messages.ConcatenateChunkMsg(chunk)
				}
			}
		}()

		i.status = Streaming
		return messages.AgentMessageAddedMsg{
			StreamChan: stream,
		}
	}
}
