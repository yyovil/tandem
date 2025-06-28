package tools

import (
	"context"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	docker "github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/yyovil/tandem/internal/utils"
)

type DockerExec ToolCall

var dockerClient *docker.Client
var kaliImage = "kali:withtools"
var containerId string

var DockerExecTool = Tool{
	Name:        DOCKER_EXEC,
	Description: "Executes a command in a Docker container.",
	Parameters: map[string]Param{
		"command": {
			Type:        utils.TypeString,
			Description: "The command to execute in the Docker container.",
		},
		"args": {
			Type: utils.TypeArray,
			Items: &Param{
				Type:        utils.TypeString,
				Description: "An argument to pass to the command.",
			},
			Description: "The arguments to pass to the command.",
		},
	},
	Required: []string{"command"},
}

func (d DockerExec) Execute(toolCallId string) ToolResponse {
	toolCallFailureResponse := ToolResponse{
		Name:       DOCKER_EXEC,
		ToolCallId: toolCallId,
		Status:     Failure,
	}
	toolCallSuccessResponse := ToolResponse{
		Name:       DOCKER_EXEC,
		ToolCallId: toolCallId,
		Status:     Success,
	}

	if dockerClient == nil {
		var err error
		dockerClient, err = docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
		if err != nil {
			toolCallFailureResponse.ToolCallResult.Error = errors.Wrap(err, "failed to create docker client")
			return toolCallFailureResponse
		}
	}

	ctx := context.Background()
	listOptions := container.ListOptions{
		All: true,
		Filters: filters.NewArgs(
			filters.KeyValuePair{
				Key:   "ancestor",
				Value: kaliImage,
			},
		),
	}

	summarySlice, err := dockerClient.ContainerList(ctx, listOptions)
	if err != nil {
		toolCallFailureResponse.ToolCallResult.Error = errors.Wrap(err, "failed to list docker containers")
		return toolCallFailureResponse
	}

	// NOTE: the whole point of going through the summary is to spot a container with a running bash exec process. so you break out soon as you found yourself one.
	for _, summary := range summarySlice {
		for _, mountPoint := range summary.Mounts {
			// NOTE: this works because we aren't going to have any other kind of mounts.
			if mountPoint.Type != mount.TypeBind {
				toolCallFailureResponse.ToolCallResult.Error = errors.New("bind mount not set up for container " + summary.ID)
				return toolCallFailureResponse
			}
		}

		// NOTE: after this switch expression, we assume that we have a container session for us.
		switch summary.State {
		case container.StateRunning:
			containerId = summary.ID
		case container.StatePaused:
			if err := dockerClient.ContainerUnpause(ctx, summary.ID); err != nil {
				toolCallFailureResponse.ToolCallResult.Error = errors.Wrap(err, "failed to unpause docker container")
				return toolCallFailureResponse
			}
			containerId = summary.ID

		case container.StateExited, container.StateCreated:
			{
				if err := dockerClient.ContainerStart(ctx, summary.ID, container.StartOptions{}); err != nil {
					toolCallFailureResponse.ToolCallResult.Error = errors.Wrap(err, "failed to start docker container")
					return toolCallFailureResponse
				}
				containerId = summary.ID
			}

		default:
			// TODO: where are the bind mounts?
			config := &container.Config{
				Image: kaliImage,
				Tty:   true,
				Cmd:   []string{"/bin/bash"},
				// continue... you were thinking about the stdin, stdout and stderr.
			}
			resp, err := dockerClient.ContainerCreate(ctx, config, nil, nil, nil, "")
			if err != nil {
				toolCallFailureResponse.ToolCallResult.Error = errors.Wrap(err, "failed to create docker container")
				return toolCallFailureResponse
			}

			if err := dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
				toolCallFailureResponse.ToolCallResult.Error = errors.Wrap(err, "failed to start docker container")
				return toolCallFailureResponse
			}

			containerId = resp.ID
		}

		if containerId != "" {
			break
		}
	}

	var commandLine []string
	for param, value := range d.Args {
		if param == "command" {
			commandLine = append(commandLine, value.(string))
		}
		if param == "args" {
			if args, ok := value.([]any); ok {
				for _, arg := range args {
					if argStr, ok := arg.(string); ok {
						commandLine = append(commandLine, argStr)
					}
				}
			}
		}
	}

	execOptions := container.ExecOptions{
		Cmd:          commandLine,
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Tty:          true,
		// NOTE: you may want to set the console size here. that would get you nice formatted output string to render in the chat view.
	}

	// NOTE: you won't be able to run tui apps then like this. because you are simply creating exec processes on fly inside the bash shell you have got running in the container.
	execCreateResponse, err := dockerClient.ContainerExecCreate(ctx, containerId, execOptions)
	if err != nil {
		toolCallFailureResponse.ToolCallResult.Error = errors.Wrap(err, "failed to create an exec process running a bash session in the docker container")
		return toolCallFailureResponse
	}

	hijackedResponse, err := dockerClient.ContainerExecAttach(ctx, execCreateResponse.ID, container.ExecAttachOptions{})
	if err != nil {
		toolCallFailureResponse.ToolCallResult.Error = errors.Wrap(err, "failed to attach to exec process in the docker container")
		return toolCallFailureResponse
	}
	defer hijackedResponse.Close()

	err = dockerClient.ContainerExecStart(ctx, execCreateResponse.ID, container.ExecStartOptions{})
	if err != nil {
		toolCallFailureResponse.ToolCallResult.Error = errors.Wrap(err, "failed to start exec process in the docker container")
		return toolCallFailureResponse
	}

	execInspect, err := dockerClient.ContainerExecInspect(ctx, execCreateResponse.ID)
	if err != nil {
		toolCallFailureResponse.ToolCallResult.Error = errors.Wrap(err, "failed to inspect exec process in the docker container")
		return toolCallFailureResponse
	}

	output, _ := io.ReadAll(hijackedResponse.Reader)

	toolCallSuccessResponse.ToolCallResult.Output = map[string]any{
		"output":    string(output),
		"exit_code": execInspect.ExitCode,
	}

	return toolCallSuccessResponse
}
