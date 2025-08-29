package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/yyovil/tandem/internal/logging"
)

const (
	TerminalToolName = "terminal"
	DockerImage      = "kali:headless"
)

type TerminalArgs struct {
	Command string   `json:"command"`
	Args    []string `json:"args,omitempty"`
}

type Terminal struct {
	hijackedResponse *types.HijackedResponse
	client           *client.Client
	init             sync.Once
	initErr          error
	containerId      string
}

var terminal *Terminal

func (term *Terminal) initialise() error {
	term.init.Do(func() {
		term.client, term.initErr = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if term.initErr != nil {
			term.initErr = fmt.Errorf("failed to create docker client: %w", term.initErr)
		}
	})

	return term.initErr
}

// Client returns the APIClient
func Client() *client.Client {
	if err := terminal.initialise(); err != nil {
		logging.Error("initialisation failed", err)
	}
	return terminal.client
}

func NewDockerCli() BaseTool {
	terminal = &Terminal{
		init:             sync.Once{},
		initErr:          nil,
		containerId:      "",
		hijackedResponse: nil,
	}

	terminal.client = Client()

	return terminal
}

func (term *Terminal) Info() ToolInfo {
	return ToolInfo{
		Name:        TerminalToolName,
		Description: "A tool to execute arbitary shell commands in a kali linux container. leave args empty for no arguments.",
		Parameters: map[string]any{
			"command": map[string]any{
				"type":        "string",
				"description": "shell command to run",
			},
			"args": map[string]any{
				"type":        "array",
				"description": "list of arguments for the command",
				"items": map[string]any{
					"type":        "string",
					"description": "argument for the command",
				},
			},
		},
		Required: []string{"command", "args"},
	}
}

func (term *Terminal) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var args TerminalArgs
	if err := json.Unmarshal([]byte(call.Input), &args); err != nil {
		return NewTextErrorResponse("Failed to parse docker cli arguments: " + err.Error()), nil
	}

	if args.Command == "" {
		return NewTextErrorResponse("command is required for DockerCli tool"), nil
	}

	// Build command slice (no trailing newline; we will exec directly)
	cmd := []string{args.Command}
	if len(args.Args) != 0 {
		cmd = append(cmd, args.Args...)
	}

	output, execErr := term.ExecuteCmd(ctx, cmd)
	if execErr != nil {
		return NewTextErrorResponse(execErr.Error()), nil
	}
	return ToolResponse{
		Type:    ToolResponseTypeText,
		Content: output,
		IsError: false,
	}, nil
}

// ExecuteCmd ensures the container is running and executes the provided command array inside it,
// returning the combined stdout/stderr output or an error.
func (term *Terminal) ExecuteCmd(ctx context.Context, cmd []string) (string, error) {
	// Ensure we have a container ID; find one by ancestor image if missing
	if term.containerId == "" {
		summaries, err := term.client.ContainerList(ctx, container.ListOptions{
			All:     true,
			Filters: filters.NewArgs(filters.Arg("ancestor", DockerImage)),
		})
		if err != nil {
			logging.Error("Failed to list containers", err)
			return "", fmt.Errorf("Failed to list containers: %w", err)
		}
		for _, summary := range summaries {
			if summary.Image == DockerImage && summary.State == container.StateRunning {
				term.containerId = summary.ID
				break
			}
			if term.containerId == "" && summary.Image == DockerImage {
				term.containerId = summary.ID
			}
		}
	}

	if term.containerId == "" {
		return "", fmt.Errorf("couldn't find a container using %s image.", DockerImage)
	}

	inspectRes, err := term.client.ContainerInspect(ctx, term.containerId)
	if err != nil {
		logging.Error(fmt.Sprintf("Failed to inspect container %s", term.containerId), err)
		return "", fmt.Errorf("Failed to inspect container: %w", err)
	}

	if !inspectRes.State.Running {
		if err := term.GetRunning(ctx, term.containerId, inspectRes.State.Status); err != nil {
			return "", fmt.Errorf("couldn't get the container: %s running", term.containerId)
		}
	}

	execResp, err := term.client.ContainerExecCreate(ctx, term.containerId, container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	})
	if err != nil {
		return "", fmt.Errorf("Failed to create exec: %w", err)
	}

	attachResp, err := term.client.ContainerExecAttach(ctx, execResp.ID, container.ExecStartOptions{})
	if err != nil {
		return "", fmt.Errorf("Failed to attach exec: %w", err)
	}
	defer attachResp.Close()

	outputBytes, err := io.ReadAll(attachResp.Reader)
	if err != nil {
		return "", fmt.Errorf("Failed to read exec output: %w", err)
	}

	output := string(outputBytes)
	logging.Debug(fmt.Sprintf("terminal exec output (%s): %s", strings.Join(cmd, " "), truncateForLog(output)))
	return output, nil
}

// NOTE: GetRunning gets a docker container to container.StateRunning.
func (term *Terminal) GetRunning(ctx context.Context, containerId string, currentState container.ContainerState) error {
	switch currentState {
	case container.StateExited, container.StateCreated:
		if err := term.client.ContainerStart(ctx, containerId, container.StartOptions{}); err != nil {
			return fmt.Errorf("failed to start container %s: %w", containerId, err)
		}
	case container.StatePaused:
		if err := term.client.ContainerUnpause(ctx, containerId); err != nil {
			return fmt.Errorf("failed to unpause container %s: %w", containerId, err)
		}
	}
	return nil
}

func truncateForLog(s string) string {
	if len(s) > 800 { // keep logs manageable
		return s[:800] + "..."
	}
	return s
}
