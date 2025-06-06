package components

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	fp "github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/key"
	vp "github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yyovil/tui/internal/utils"
)

type FilePicker struct {
	showFilePicker bool
	viewport       vp.Model
	filepicker     fp.Model
	width          int
	height         int
	selectedFiles  []string //slice storing the path for the selected files.
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
				// this is to set the scroll position to top when you select a dir.
				fpc.viewport.GotoTop()
				fpc.viewport.Update(msg)
			} else {
				fpc.selectedFiles = append(fpc.selectedFiles, fpc.filepicker.FileSelected)
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
		MaxHeight(1).
		Padding(0, 1, 0, 1).
		AlignVertical(lipgloss.Center).
		Border(lipgloss.InnerHalfBlockBorder(), false, true).
		Background(lipgloss.Color("#343a40"))

	fpSelectedStyle := fpc.filepicker.Styles.Selected
	if len(fpc.selectedFiles) == 0 {
		s.WriteString("Pick a file")
	} else if len(fpc.selectedFiles) == 1 {
		footerStyle = footerStyle.BorderForeground(lipgloss.Color("212"))
		s.WriteString("Selected file: " + fpSelectedStyle.Render(fpc.selectedFiles[0]))
	} else {
		footerStyle = footerStyle.BorderForeground(lipgloss.Color("212"))
		s.WriteString("Total files selected: " + fpSelectedStyle.Render(fmt.Sprintf("%d", len(fpc.selectedFiles))))
	}

	return footerStyle.Render(s.String())
}

func (fpc FilePicker) GetSelectedFile() (utils.Attachment, error) {
	fs := fpc.filepicker.FileSelected
	if fs != "" {
		// fileStat, err := os.Stat(fs)

		// if err != nil {
		// 	return utils.Attachment{
		// 		Filepath: fpc.filepicker.FileSelected,
		// 		Content:  []byte{},
		// 	}, err
		// }

		// if fileStat.IsDir() {
		// FEATURE: in future we would like to support uploading multiple files at 1 level depth by selecting a dir.
		// return utils.Attachment{

		// 		Content: []byte{},
		// 	}, errors.New("can't get the selected dir")
		// }

		content, err := os.ReadFile(fs)
		if err != nil {
			log.Println("error reading file:", err.Error())
			return utils.Attachment{
				Filepath: "",
				Url:      "",
				Content:  "",
			}, err
		}

		return utils.Attachment{
			Filepath: fpc.filepicker.FileSelected,
			Url:      "",
			MimeType: strings.Split(http.DetectContentType(content), ";")[0],
			Content:  string(content),
		}, nil

	} else {
		return utils.Attachment{
			Filepath: "",
			Content:  "",
		}, nil
	}
}

func NewFilePicker() FilePicker {
	return FilePicker{
		showFilePicker: false,
		filepicker:     fp.New(),
	}
}

type AttachmentMsg struct {
	Attachment utils.Attachment
}

/*
TODO: support multifile selection.
BUG: we need to know the file object at the cursor position when Open action is performed.
*/
