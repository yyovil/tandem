package components

import (
	// "bufio"
	// "encoding/json"
	"fmt"
	// "log"
	// "net/http"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	vp "github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yyovil/tandem/internal/components/messages"
	"github.com/yyovil/tandem/internal/layout"
	// "github.com/yyovil/tandem/internal/utils"

)

// type Status string

// const (
// 	Requesting    Status = "Requesting"
// 	Streaming     Status = "Streaming"
// 	ToolCall      Status = "Tool call"
// 	ToolCompleted Status = "Tool completed"
// 	Idle          Status = "Idle"
// )

type Input struct {
	// status        Status
	stream        chan tea.Msg
	width, height int
	UserPrompt    string
	spinner       spinner.Model
	textarea      textarea.Model
	FilePicker    FilePicker
	// TODO: out this and put in a dedicated layout file.
	leftpane, rightpane vp.Model

	leftPaneMessages []tea.Msg
}

type InputKeyMap struct {
	ShowFilePicker,
	// Send,
	Quit,
	PageDown,
	PageUp,
	HalfPageUp,
	HalfPageDown key.Binding
	// TODO: add a keybinding for toggling the tool execution view.
}

var inputKeyMap = InputKeyMap{
	ShowFilePicker: key.NewBinding(
		key.WithKeys("ctrl+o"),
		key.WithHelp("ctrl+o", "attach file"),
	),
	// Send: key.NewBinding(
	// 	key.WithKeys("enter"),
	// 	key.WithHelp("enter", "send message"),
	// ),
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
	i.textarea.ShowLineNumbers = false
	return tea.Batch(
		textarea.Blink,
		i.spinner.Tick,
		i.textarea.Focus(),
	)
}

func (i *Input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
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
			i.textarea.Blur()

		// case key.Matches(msg, inputKeyMap.Send):
		// 	if !i.FilePicker.showFilePicker {

		// 		if i.textarea.Value() == "" || i.status != Idle {
		// 			return i, nil
		// 		}

		// 		i.userPrompt = i.textarea.Value()
		// 		cmds = append(cmds, i.sendRunRequestCmd(), messages.AddUserMessageCmd(i.userPrompt, i.FilePicker.selectedFiles))
		// 		i.leftpane.GotoBottom()

		// 		i.textarea.Reset()
		// 		i.FilePicker.viewport.GotoTop()
		// 		i.FilePicker.filepicker.FileSelected = ""

		// 		return i, tea.Batch(cmds...)
		// 	}

		// 	_, cmd = i.FilePicker.Update(msg)
		// 	cmds = append(cmds, cmd)
		// 	return i, tea.Batch(cmds...)

		case key.Matches(msg, inputKeyMap.Quit):
			if !i.textarea.Focused() {
				i.textarea.Focus()
			}

			if !i.FilePicker.showFilePicker {
				return i, tea.Quit
			}
		default:
			if !i.FilePicker.showFilePicker {
				cmd = i.textarea.Focus()
				cmds = append(cmds, cmd)
				i.textarea, cmd = i.textarea.Update(msg)
				cmds = append(cmds, cmd)
				return i, tea.Batch(cmds...)
			}
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

	case messages.RunStartedMsg:
		// blocking call to receive the first chunk of the stream
		// i.status = Streaming
		agentMessage := messages.AgentMessage{
			Width:   i.leftpane.Width,
			Content: &strings.Builder{},
		}
		firstChunk, _ := <-i.stream
		if v, ok := firstChunk.(messages.RunResponseContentMsg); ok {
			agentMessage.Content.WriteString(v.Content)
			i.leftPaneMessages = append(i.leftPaneMessages, agentMessage)
			i.leftpane.GotoBottom()

			return i, messages.ListenOnStreamChanCmd(i.stream)
		}

	case messages.RunResponseContentMsg:
		// i.status = Streaming
		// Update the last message in place instead of appending a new one
		if len(i.leftPaneMessages) > 0 {
			lastMsgIndex := len(i.leftPaneMessages) - 1
			lastMsg := i.leftPaneMessages[lastMsgIndex]

			if agentMsg, ok := lastMsg.(messages.AgentMessage); ok {
				agentMsg.Content.WriteString(msg.Content)
				i.leftPaneMessages[lastMsgIndex] = agentMsg
				i.leftpane.GotoBottom()

			} else {
				agentMessage := messages.AgentMessage{
					Width:   i.leftpane.Width,
					Content: &strings.Builder{},
				}
				agentMessage.Content.WriteString(msg.Content)
				i.leftPaneMessages = append(i.leftPaneMessages, agentMessage)
				i.leftpane.GotoBottom()
			}
		}
		return i, messages.ListenOnStreamChanCmd(i.stream)

	case messages.ToolCallStartedMsg:
		// i.status = ToolCall
		toolCallMsg := messages.ToolExecutionMessage{
			Width:        i.leftpane.Width,
			Event:        msg.Event,
			ToolCallName: msg.Tool.ToolName,
		}

		i.leftPaneMessages = append(i.leftPaneMessages, toolCallMsg)
		i.leftpane.GotoBottom()
		return i, messages.ListenOnStreamChanCmd(i.stream)

	case messages.ToolCallCompletedMsg:
		// i.status = ToolCompleted
		if len(i.leftPaneMessages) > 0 {
			lastMsgIndex := len(i.leftPaneMessages) - 1
			lastMsg := i.leftPaneMessages[lastMsgIndex]

			if teMsg, ok := lastMsg.(messages.ToolExecutionMessage); ok {
				if msg.Tool.Result != "" {
					teMsg.ToolCallResult = msg.Tool.Result
				}

				if msg.Content != "" {
					teMsg.Content = msg.Content
				}

				i.leftPaneMessages[lastMsgIndex] = teMsg
				i.leftpane.GotoBottom()

				return i, messages.ListenOnStreamChanCmd(i.stream)
			}
		}

	case messages.RunResponseCompletedMsg:
		// i.status = Idle
		return i, nil
	}

	i.spinner, cmd = i.spinner.Update(msg)
	cmds = append(cmds, cmd)
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

	// content := ""
	var content strings.Builder
	for _, msg := range i.leftPaneMessages {
		switch m := msg.(type) {
		case messages.UserMessage:
			content.WriteString(m.View() + "\n")
		case messages.AgentMessage:
			content.WriteString(m.View() + "\n")
		case messages.ToolExecutionMessage:
			content.WriteString(m.View() + "\n")
		}
	}

	i.leftpane.SetContent(content.String())

	panes := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPaneStyle.Render(i.leftpane.View()),
		rightPaneStyle.Render(i.rightpane.View()),
	)

	if i.FilePicker.showFilePicker {
		x := (i.width - i.FilePicker.viewport.Width) / 2
		y := (i.height - i.FilePicker.viewport.Height) / 2

		fg := i.FilePicker.View()
		bg := lipgloss.JoinVertical(lipgloss.Top, panes, inputStyle.Render(i.textarea.View())+"\n"+i.footerView())

		return layout.Composite(x, y, fg, bg)
	}

	// if i.status != Idle {
	// 	return lipgloss.JoinVertical(lipgloss.Top,
	// 		panes,
	// 		i.headerView()+"\n"+inputStyle.Render(i.textarea.View())+"\n"+i.footerView())
	// }

	return lipgloss.JoinVertical(lipgloss.Top,
		panes,
		inputStyle.Render(i.textarea.View())+"\n"+i.footerView())
}

// displays the status of the run request to the agent.
// func (i Input) headerView() string {
// 	spinnerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffafcc")).Italic(true).PaddingLeft(1)
// 	return spinnerStyle.Render(string(i.status) + i.spinner.View())
// }

// displays the attachement
func (i Input) footerView() string {
	var s strings.Builder
	footerStyle := lipgloss.
		NewStyle().
		Height(1).
		Width(i.width-2).
		MaxHeight(2).
		Border(lipgloss.InnerHalfBlockBorder(), false).
		BorderLeft(true).
		BorderRight(true).
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

	// statusStyle := lipgloss.
	// 	NewStyle().
	// 	Width(10).
	// 	Height(1).
	// 	Border(lipgloss.InnerHalfBlockBorder(), false).
	// 	BorderRight(true).
	// 	Background(lipgloss.Color("#fb5607")).
	// 	AlignHorizontal(lipgloss.Center)

	// return lipgloss.JoinHorizontal(lipgloss.Left, footerStyle.Render(s.String()), statusStyle.Render(string(i.status)))
	return footerStyle.Render(s.String())
}

func NewInput() Input {
	return Input{
		spinner:          spinner.New(spinner.WithSpinner(spinner.Meter)),
		// status:           Idle,
		leftPaneMessages: []tea.Msg{},
		width:            0,
		height:           0,
		UserPrompt:       "",
		textarea:         textarea.New(),
		FilePicker:       NewFilePicker(),
		leftpane:         vp.New(0, 0),
		rightpane:        vp.New(0, 0),
	}
}

// func (i *Input) sendRunRequestCmd() tea.Cmd {
// 	// i.status = Requesting

// 	return func() tea.Msg {
// 		attachments, err := i.FilePicker.GetSelectedFiles()
// 		if err != nil {
// 			// i.status = Idle
// 			log.Println("error getting the attachment:", err.Error())
// 			return func() tea.Msg {
// 				return nil
// 			}
// 		}

// 		req, err := utils.GetPostRequest(i.UserPrompt, attachments)
// 		if err != nil {
// 			// i.status = Idle
// 			log.Println("error creating request:", err.Error())
// 			return func() tea.Msg {
// 				return nil
// 			}
// 		}

// 		// UNBUFFERED CHANNEL
// 		stream := make(chan tea.Msg)

// 		go func() {
// 			defer close(stream)

// 			client := &http.Client{}
// 			resp, err := client.Do(req)
// 			if err != nil {
// 				i.status = Idle

// 				log.Println("error sending request:", err.Error())
// 				// TODO: show a user feedback for this error
// 				return
// 			}
// 			defer resp.Body.Close()

// 			if resp.StatusCode != http.StatusOK {
// 				i.status = Idle
// 				log.Println("error: received non-200 response:", resp.Status)
// 				// TODO: show a user feedback for this error
// 				return
// 			}

// 			scanner := bufio.NewScanner(resp.Body)
// 			scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
// 				if atEOF && len(data) == 0 {
// 					return 0, nil, nil
// 				}
// 				if idx := strings.Index(string(data), "\n\n"); idx >= 0 {
// 					return idx + 2, data[:idx], nil
// 				}
// 				if atEOF {
// 					return len(data), data, nil
// 				}
// 				return 0, nil, nil
// 			})

// 			for scanner.Scan() {
// 				chunk := scanner.Bytes()

// 				var rr messages.RunResponse
// 				if err := json.Unmarshal(chunk, &rr); err != nil {
// 					log.Println("error decoding chunk:", err.Error())
// 				} else {
// 					switch rr.Event {

// 					case messages.RunResponseContent:

// 						stream <- messages.RunResponseContentMsg(rr)

// 					case messages.ToolCallStarted:
// 						stream <- messages.ToolCallStartedMsg(rr)

// 					case messages.ToolCallCompleted:
// 						stream <- messages.ToolCallCompletedMsg(rr)

// 					default:
// 						i.status = Idle
// 						log.Println("unknown event type:", rr)
// 					}
// 				}
// 			}
// 		}()

// 		i.status = Streaming
// 		i.stream = stream
// 		return messages.RunStartedMsg{}
// 	}
// }
