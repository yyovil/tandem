package bubbles

// import (
// 	"log"
// 	"net/http"
// 	"os"
// 	"strings"

// 	fp "github.com/charmbracelet/bubbles/filepicker"
// 	"github.com/charmbracelet/bubbles/key"
// 	vp "github.com/charmbracelet/bubbles/viewport"
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"
// 	"github.com/yyovil/tandem/internal/utils"
// )

// type FilePicker struct {
// 	showFilePicker bool
// 	viewport       vp.Model
// 	filepicker     fp.Model
// 	width          int
// 	height         int
// 	SelectedFiles  []string //slice storing the path for the selected files.
// }

// type FilePickerKeyMap struct {
// 	Cancel key.Binding
// }

// var filePickerKeyMap = FilePickerKeyMap{
// 	Cancel: key.NewBinding(
// 		key.WithKeys("esc"),
// 		key.WithHelp("esc", "cancel"),
// 	),
// }

// func (fpc *FilePicker) Init() tea.Cmd {

// 	if pwd, err := os.UserHomeDir(); err != nil {
// 		log.Println("$HOME not set", err.Error())
// 		os.Exit(1)
// 	} else {
// 		fpc.filepicker.CurrentDirectory = pwd
// 		fpc.filepicker.ShowSize = true
// 		fpc.filepicker.DirAllowed = true
// 		fpc.filepicker.FileAllowed = true
// 	}
// 	return fpc.filepicker.Init()
// }

// func (fpc *FilePicker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	var (
// 		cmds []tea.Cmd
// 		cmd  tea.Cmd
// 	)

// 	switch msg := msg.(type) {
// 	case tea.WindowSizeMsg:
// 		fpc.width = (msg.Width * 80) / 100
// 		fpc.height = msg.Height / 2

// 		fpc.viewport = vp.New(fpc.width, fpc.height)

// 		fpc.viewport.GotoTop()
// 		fpc.filepicker.SetHeight(fpc.height)

// 	case tea.KeyMsg:
// 		switch {
// 		case key.Matches(msg, filePickerKeyMap.Cancel):
// 			fpc.showFilePicker = false
// 			return fpc, nil

// 		case key.Matches(msg, fpc.filepicker.KeyMap.Open):
// 			if fpc.filepicker.FileSelected != "" {
// 				if fileStat, err := os.Stat(fpc.filepicker.FileSelected); err != nil {
// 					log.Println("Error getting file stat", err.Error())
// 					// TODO: show status to user.
// 				} else if fileStat.IsDir() {
// 					fpc.viewport.GotoTop()
// 					fpc.viewport.Update(msg)
// 				}
// 			}
// 		}
// 	}

// 	fpc.filepicker, cmd = fpc.filepicker.Update(msg)
// 	cmds = append(cmds, cmd)

// 	if didSelect, path := fpc.filepicker.DidSelectFile(msg); didSelect {
// 		fpc.filepicker.FileSelected = path
// 		if fpc.filepicker.FileSelected != "" {
// 			if fileStat, err := os.Stat(fpc.filepicker.FileSelected); err != nil {
// 				log.Println("Error getting file stat", err.Error())
// 				// TODO: show status to user.
// 			} else if fileStat.IsDir() {
// 				// this is to set the scroll position to top when you select a dir.
// 				fpc.viewport.GotoTop()
// 				fpc.viewport.Update(msg)
// 			} else {
// 				fpc.SelectedFiles = append(fpc.SelectedFiles, fpc.filepicker.FileSelected)
// 			}
// 		}
// 		// log.Println("fpc.SelectedFiles: ", fpc.SelectedFiles)
// 	}

// 	fpc.viewport, cmd = fpc.viewport.Update(msg)
// 	cmds = append(cmds, cmd)

// 	return fpc, tea.Batch(cmds...)
// }

// func (fpc *FilePicker) View() string {

// 	if !fpc.showFilePicker {
// 		return ""
// 	}

// 	var s strings.Builder
// 	s.WriteString(fpc.filepicker.View())
// 	fpc.viewport.SetContent(s.String())
// 	fpc.viewport.Style = lipgloss.NewStyle().
// 		// Width(fpc.width).
// 		// Height(fpc.height).
// 		// MaxWidth(fpc.width).
// 		// MaxHeight(fpc.height).
// 		Border(lipgloss.NormalBorder())

// 	return fpc.viewport.View()
// }


// func (fpc *FilePicker) GetSelectedFiles() (attachments []utils.Attachment, error error) {
// 	if len(fpc.SelectedFiles) > 0 {
// 		// fileStat, err := os.Stat(fs)

// 		// if err != nil {
// 		// 	return utils.Attachment{
// 		// 		Filepath: fpc.filepicker.FileSelected,
// 		// 		Content:  []byte{},
// 		// 	}, err
// 		// }

// 		// if fileStat.IsDir() {
// 		// FEATURE: in future we would like to support uploading multiple files at 1 level depth by selecting a dir.
// 		// return utils.Attachment{

// 		// 		Content: []byte{},
// 		// 	}, errors.New("can't get the selected dir")
// 		// }

// 		for _, selectedFile := range fpc.SelectedFiles {
// 			content, err := os.ReadFile(selectedFile)
// 			if err != nil {
// 				log.Println("error reading a file:", err.Error())
// 				return nil, err
// 			} else {
// 				attachments = append(attachments, utils.Attachment{
// 					Filepath: selectedFile,
// 					MimeType: strings.Split(http.DetectContentType(content), ";")[0],
// 					Content:  string(content),
// 				})
// 			}
// 		}
// 		fpc.SelectedFiles = nil

// 		return attachments, nil
// 	} else {
// 		return nil, nil
// 	}
// }

// func NewFilePicker() FilePicker {
// 	return FilePicker{
// 		showFilePicker: false,
// 		filepicker:     fp.New(),
// 	}
// }

// type AttachmentMsg struct {
// 	Attachment utils.Attachment
// }

// /*
// BUG: we need to know the file object at the cursor position when Open action is performed.
// */
