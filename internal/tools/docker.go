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
	"github.com/yaydraco/tandem/internal/logging"
)

/*
TODO: implement docker_cli adhering to the BaseTool interface.
*/

const (
	DockerCliToolName = "docker_cli"
	DockerImage       = "kali:withtools"
)

type DockerCliArgs struct {
	Command string   `json:"command"`
	Args    []string `json:"args,omitempty"`
}

type DockerCli struct {
	hijackedResponse *types.HijackedResponse
	client           *client.Client
	init             sync.Once
	initErr          error
	containerId      string
}

var dockerCli *DockerCli

func (cli *DockerCli) initialise() error {
	cli.init.Do(func() {
		cli.client, cli.initErr = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if cli.initErr != nil {
			cli.initErr = fmt.Errorf("failed to create docker client: %w", cli.initErr)
		}
	})

	return cli.initErr
}

// Client returns the APIClient
func Client() *client.Client {
	if err := dockerCli.initialise(); err != nil {
		logging.Error("initialisation failed", err)
	}
	return dockerCli.client
}

// NOTE: when this tool is called, its expected that it was during the initialisation time,
func NewDockerCli() BaseTool {
	dockerCli = &DockerCli{
		init:             sync.Once{},
		initErr:          nil,
		containerId:      "",
		hijackedResponse: nil,
	}
	dockerCli.client = Client()

	return dockerCli
}

func (cli *DockerCli) Info() ToolInfo {
	return ToolInfo{
		Name:        DockerCliToolName,
		Description: "A tool to execute arbitary shell commands in a docker container.",
		Parameters: map[string]any{
			"command": map[string]any{
				"type":        "string",
				"description": "shell command to run",
			},
			"args": map[string]any{
				"type":        "array",
				"description": "arguments for the command",
				"items": map[string]any{
					"type":        "string",
					"description": "argument for the command",
				},
			},
		},
		Required: []string{"command"},
	}
}

func (cli *DockerCli) Run(ctx context.Context, call ToolCall) (ToolResponse, error) {
	var args DockerCliArgs
	if err := json.Unmarshal([]byte(call.Input), &args); err != nil {
		return NewTextErrorResponse("Failed to parse docker cli arguments: " + err.Error()), nil
	}

	if args.Command == "" {
		return NewTextErrorResponse("command is required for DockerCli tool"), nil
	}

	commandLine := args.Command
	if len(args.Args) != 0 {
		commandLine += " " + strings.Join(args.Args, " ")
	}
	// NOTE: when using tty in normal mode, input is line buffered, implying you need to press enter to send the command. thus appending a newline.
	commandLine += "\n"

	if cli.containerId == "" {
		summaries, err := cli.client.ContainerList(ctx, container.ListOptions{
			All:     true,
			Filters: filters.NewArgs(filters.Arg("ancestor", DockerImage)),
		})
		if err != nil {
			return NewTextErrorResponse("Failed to list containers: " + err.Error()), nil
		}

		for _, summary := range summaries {
			if summary.Image == DockerImage && summary.State == container.StateRunning {
				cli.containerId = summary.ID
				break
			}

			if summary.Image == DockerImage {
				cli.containerId = summary.ID
				break
			}
		}
	}

	// NOTE: we are not creating a container if not found in the summaries because it should be created during the installation.
	if cli.containerId == "" {
		return NewTextErrorResponse(fmt.Sprintf("couldn't find a container using %s image.", DockerImage)), nil
	}

	inspectRes, err := cli.client.ContainerInspect(ctx, cli.containerId)
	if err != nil {
		return NewTextErrorResponse("Failed to inspect container: " + err.Error()), nil
	}

	if !inspectRes.State.Running {
		if err := cli.GetRunning(ctx, cli.containerId, inspectRes.State.Status); err != nil {
			return NewTextErrorResponse(fmt.Sprintf("couldn't get the container: %s running.", cli.containerId)), nil
		}
	}

	if cli.hijackedResponse == nil {
		hijackedResp, err := cli.client.ContainerAttach(ctx, cli.containerId, container.AttachOptions{
			Stream: true,
			Stdin:  true,
			Stdout: true,
			Stderr: true,
		})
		if err != nil {
			return NewTextErrorResponse("Failed to attach to container: " + err.Error()), nil
		}
		defer hijackedResp.Close()
		cli.hijackedResponse = &hijackedResp
	}

	if _, err := cli.hijackedResponse.Conn.Write([]byte(commandLine)); err != nil {
		return NewTextErrorResponse("Failed to write command to container: " + err.Error()), nil
	}

	var output []byte
	for {
		_, err := cli.hijackedResponse.Reader.Read(output)
		if err != nil {
			if err == io.EOF {
				break // End of stream
			}
			return NewTextErrorResponse("Failed to read from container: " + err.Error()), nil
		}
	}

	return ToolResponse{
		Type:    ToolResponseTypeText,
		Content: string(output),
		IsError: false,
	}, nil
}

// NOTE: GetRunning gets a docker container to container.StateRunning.
func (cli *DockerCli) GetRunning(ctx context.Context, containerId string, currentState container.ContainerState) error {
	switch currentState {
	case container.StateExited, container.StateCreated:
		if err := cli.client.ContainerStart(ctx, containerId, container.StartOptions{}); err != nil {
			return fmt.Errorf("failed to start container %s: %w", containerId, err)
		}
	case container.StatePaused:
		if err := cli.client.ContainerUnpause(ctx, containerId); err != nil {
			return fmt.Errorf("failed to unpause container %s: %w", containerId, err)
		}
	}
	return nil
}
