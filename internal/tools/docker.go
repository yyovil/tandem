package tools

/*
TODO: implement docker_exec adhering to the BaseTool interface.
*/

const (
	ShellToolName = "shell"
)

type ShellParams struct {
	Command string `json:"command"`
}
