package components

import (
	"bufio"
	"log"
	"strings"

	"net/http"

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
	Streaming  Status = "streaming"
	Requesting Status = "requesting"
	Idle       Status = "idle"
)

type Input struct {
	//TODO: create a status indicator above the input cmp.
	status Status

	width, height int
	userPrompt    string
	textarea      textarea.Model
	FilePicker    FilePicker
	// TODO: out this and put in a dedicated layout file.
	leftpane, rightpane vp.Model

	leftPaneMessages []any //TODO: this out, use a better type
}

type InputKeyMap struct {
	ShowFilePicker, Send, Quit key.Binding
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
		case key.Matches(msg, inputKeyMap.ShowFilePicker):
			cmd = i.FilePicker.Init()
			cmds = append(cmds, cmd)
			i.FilePicker.showFilePicker = true
		case key.Matches(msg, inputKeyMap.Send):
			// TODO: only send the message if the i.status == Idle
			if !i.FilePicker.showFilePicker {

				if i.textarea.Value() == "" {
					return i, nil
				}

				i.userPrompt = i.textarea.Value()
				attachmentName := i.FilePicker.filepicker.FileSelected
				// Send a command to add the user message
				cmds = append(cmds, messages.AddUserMessageCmd(i.userPrompt, attachmentName), sendRunRequestCmd(i.userPrompt))

				i.textarea.Reset()
				i.FilePicker.viewport.GotoTop()
				i.FilePicker.filepicker.FileSelected = ""
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
		i.leftpane, cmd = i.leftpane.Update(msg)
		cmds = append(cmds, cmd)

		i.rightpane.Width = rightPaneWidth
		i.rightpane.Height = paneHeight
		i.rightpane, cmd = i.rightpane.Update(msg)
		cmds = append(cmds, cmd)

	case messages.UserMessageAddedMsg:
		i.leftPaneMessages = append(i.leftPaneMessages, msg.UserMessage)

		// Build the left pane content from all messages
		var leftContent strings.Builder
		for _, um := range i.leftPaneMessages {
			umsg, ok := um.(messages.UserMessage)
			if ok {

				umsg.Width = i.leftpane.Width
				umsg.Height = i.leftpane.Height

				leftContent.WriteString(umsg.View())
				leftContent.WriteString("\n\n")
			}
		}

		i.leftpane.SetContent(leftContent.String())
	case messages.AgentMessageAddedMsg:
		i.status = Streaming
		agentMessage := messages.AgentMessage{
			StreamChan: msg.StreamChan,
			Content:    "",
		}

		i.leftPaneMessages = append(i.leftPaneMessages, agentMessage)

		var agentResponse strings.Builder

		for _, um := range i.leftPaneMessages {
			amsg, ok := um.(messages.AgentMessage)
			if ok {
				amsg.Width = i.leftpane.Width
				amsg.Height = i.leftpane.Height

				agentResponse.WriteString(amsg.View())
				agentResponse.WriteString("\n\n")
			}
		}

		log.Println("agent response: ", agentResponse.String())
		i.leftpane.SetContent(agentResponse.String())
		agentMessage.Update(msg)

	case messages.EndStream:
		i.status = Idle
		return i, nil

	}

	_, cmd = i.FilePicker.Update(msg)
	cmds = append(cmds, cmd)

	i.textarea, cmd = i.textarea.Update(msg)
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
		Border(lipgloss.InnerHalfBlockBorder(), false, false, false, true).
		BorderLeftForeground(lipgloss.Color("#e2a3c7")).
		// Background(lipgloss.Color("#e2a3c7")).
		Padding(1)

	panes := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPaneStyle.Render(i.leftpane.View()),
		rightPaneStyle.Render(i.rightpane.View()),
	)

	return lipgloss.JoinVertical(lipgloss.Top, panes, inputStyle.Render(i.textarea.View())+"\n"+i.footerView())
}

// displays the attachment selected.
func (i Input) footerView() string {
	var s strings.Builder
	footerStyle := lipgloss.
		NewStyle().
		Width(i.width-2).
		MaxWidth(i.width).
		Height(1).
		MaxHeight(2).
		Border(lipgloss.InnerHalfBlockBorder(), false).
		BorderLeft(true).
		BorderRight(true).
		PaddingLeft(1).
		Background(lipgloss.Color("#343a40")).
		MarginBottom(1)

	if i.FilePicker.filepicker.FileSelected == "" {
		s.WriteString("No Attachments")
	} else {
		footerStyle = footerStyle.BorderForeground(lipgloss.Color("212"))
		s.WriteString("Attachment: " + i.FilePicker.filepicker.Styles.Selected.Render(i.FilePicker.filepicker.FileSelected))
	}

	// TODO: status should appear to the far right.
	statusStyle := lipgloss.
		NewStyle().
		AlignHorizontal(lipgloss.Right)
	s.WriteString(statusStyle.Render(string(i.status)))

	return footerStyle.Render(s.String())
}

func NewInput() Input {
	return Input{
		status:           Idle,
		leftPaneMessages: []any{},
		width:            0,
		height:           0,
		userPrompt:       "",
		textarea:         textarea.New(),
		FilePicker:       NewFilePicker(),
		leftpane:         vp.New(0, 0),
		rightpane:        vp.New(0, 0),
	}
}

// requests the agent api for the agent message.
func sendRunRequestCmd(Prompt string) tea.Cmd {

	req, err := utils.Post(Prompt)
	if err != nil {
		return func() tea.Msg {
			return nil
		}
	}

	return func() tea.Msg {
		// channel to stream the text chunks.
		stream := make(chan tea.Msg)

		go func() {
			defer close(stream)
			client := &http.Client{}
			resp, err := client.Do(req)

			if err != nil {
				log.Println("error sending request:", err.Error())
			}
			
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				log.Println("error: received non-200 response:", resp.Status)
				// TODO: show a user feedback for this error
			}

			scanner := bufio.NewScanner(resp.Body)

			for scanner.Scan() {
				chunk := scanner.Text()
				stream <- messages.ConcatenateChunkMsg(chunk)
				log.Println("streaming chunk: ", chunk)
			}

			if err := scanner.Err(); err != nil {
				log.Println("error reading response body:", err.Error())
				return
			}
			log.Println("ending stream")

			stream <- messages.EndStream{}
		}()

		return messages.AgentMessageAddedMsg{
			StreamChan: stream,
		}
	}
}

/*
TODO: we don't really need 2 seperate UserMessageAddedMsg and AgentMessageAddedMsg msgs to update the leftpane viewport. Use a generic one instead.
*/
