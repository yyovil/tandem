package components

import (
	"errors"
	"log"
	"os"
	"strings"

	fp "github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/key"
	vp "github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Attachment struct {
	Name    string
	Content []byte
}

type FilePicker struct {
	showFilePicker bool
	viewport       vp.Model
	filepicker     fp.Model
	width          int
	height         int
}

type FilePickerKeyMap struct {
	Cancel key.Binding
}

var filePickerKeyMap = FilePickerKeyMap{
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
}

func (fpc *FilePicker) Init() tea.Cmd {

	if pwd, err := os.UserHomeDir(); err != nil {
		log.Println("$HOME not set", err.Error())
		os.Exit(1)
	} else {
		fpc.filepicker.CurrentDirectory = pwd
		fpc.filepicker.ShowSize = true
		fpc.filepicker.DirAllowed = true
		fpc.filepicker.FileAllowed = true
	}
	return fpc.filepicker.Init()
}

func (fpc *FilePicker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		fpc.width = msg.Width
		fpc.height = msg.Height
		fpc.viewport = vp.New((msg.Width*80)/100, msg.Height/2)
		fpc.viewport.GotoTop()
		fpc.filepicker.SetHeight(fpc.viewport.Height)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, filePickerKeyMap.Cancel):
			fpc.showFilePicker = false
			return fpc, nil

		case key.Matches(msg, fpc.filepicker.KeyMap.Open):
			if fpc.filepicker.FileSelected != "" {
				if fileStat, err := os.Stat(fpc.filepicker.FileSelected); err != nil {
					log.Println("Error getting file stat", err.Error())
					// TODO: show status to user.
				} else if fileStat.IsDir() {
					fpc.viewport.GotoTop()
					fpc.viewport.Update(msg)
				}
			}
		}
	}

	fpc.filepicker, cmd = fpc.filepicker.Update(msg)
	cmds = append(cmds, cmd)

	if didSelect, path := fpc.filepicker.DidSelectFile(msg); didSelect {
		fpc.filepicker.FileSelected = path
		if fpc.filepicker.FileSelected != "" {
			if fileStat, err := os.Stat(fpc.filepicker.FileSelected); err != nil {
				log.Println("Error getting file stat", err.Error())
				// TODO: show status to user.
			} else if fileStat.IsDir() {
				fpc.viewport.SetYOffset(0)
				fpc.viewport.Update(msg)
			}
		}
	}

	fpc.viewport, cmd = fpc.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return fpc, tea.Batch(cmds...)
}

func (fpc *FilePicker) View() string {

	if !fpc.showFilePicker {
		return ""
	}

	var s strings.Builder
	s.WriteString(fpc.filepicker.View())
	fpc.viewport.SetContent(s.String())
	fpc.viewport.Style = lipgloss.NewStyle().Border(lipgloss.NormalBorder())

	return lipgloss.Place(fpc.width, fpc.height, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Left, fpc.viewport.View(), fpc.footerView()))
}

func (fpc FilePicker) footerView() string {
	var s strings.Builder
	footerStyle := lipgloss.
		NewStyle().
		Width(fpc.viewport.Width-2).
		MaxWidth(fpc.viewport.Width).
		Height(1).
		MaxHeight(2).
		Padding(0, 1, 0, 1).
		AlignVertical(lipgloss.Center).
		Border(lipgloss.InnerHalfBlockBorder(), false, true).
		Background(lipgloss.Color("#343a40"))

	if fpc.filepicker.FileSelected == "" {
		s.WriteString("Pick a file\n\n")
	} else {
		footerStyle = footerStyle.BorderForeground(lipgloss.Color("212"))
		s.WriteString("Selected file: " + fpc.filepicker.Styles.Selected.Render(fpc.filepicker.FileSelected) + "\n")

	}
	return footerStyle.Render(s.String())
}

func (fpc FilePicker) GetSelectedFile() (Attachment, error) {
	fs := fpc.filepicker.FileSelected
	if fs != "" {
		fileStat, err := os.Stat(fs)
		if err != nil {
			return Attachment{
				Name:    fs,
				Content: []byte{},
			}, err
		}

		if fileStat.IsDir() {
			// FEATURE: in future we would like to support uploading multiple files at 1 level depth by selecting a dir.
			return Attachment{
				Name:    fs,
				Content: []byte{},
			}, errors.New("can't get the selected dir")
		}

		content, err := os.ReadFile(fs)

		if err != nil {
			log.Println("error reading file:", err.Error())
			return Attachment{
				Name:    fs,
				Content: []byte{},
			}, err
		}

		return Attachment{
			Name:    fpc.filepicker.FileSelected,
			Content: content,
		}, nil

	} else {
		return Attachment{
			Name:    fs,
			Content: []byte{},
		}, errors.New("can't get a file if not selected")
	}
}

func NewFilePicker() FilePicker {
	return FilePicker{
		showFilePicker: false,
		filepicker:     fp.New(),
	}
}

// TODO: this cmd should trigger showing the file picker component.
func AttachCmd(i Input) tea.Cmd {
	attachment, err := i.FilePicker.GetSelectedFile()
	var attachmentMsg AttachmentMsg
	if err != nil {
		// TODO: find better handling of error.
		log.Println("error getting selected file:", err.Error())
		attachmentMsg = AttachmentMsg{
			Attachment: Attachment{
				attachment.Name,
				attachment.Content,
			},
		}
	} else {
		attachmentMsg = AttachmentMsg{
			Attachment: attachment,
		}
	}
	return func() tea.Msg {
		return attachmentMsg
	}
}

type AttachmentMsg struct {
	Attachment Attachment
}

/*
TODO: support multifile selection.
BUG: we need to know the file object at the cursor position when Open action is performed.
*/
