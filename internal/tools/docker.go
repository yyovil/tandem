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

func (d DockerExec) NewClient() (*docker.Client, error) {
	client, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create docker client")
	}

	return client, nil
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
		if dockerClient, err = d.NewClient(); err != nil {
			toolCallFailureResponse.ToolCallResult.Error = err
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

	for _, summary := range summarySlice {
		for _, mountPoint := range summary.Mounts {
			// NOTE: this works because we aren't going to have any other kind of mounts.
			if mountPoint.Type != mount.TypeBind {
				toolCallFailureResponse.ToolCallResult.Error = errors.New("bind mount not set up for container " + summary.ID)
				return toolCallFailureResponse
			}
		}

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
				Image:        kaliImage,
				Tty:          true,
				Cmd:          []string{"/bin/bash"},
				AttachStdin:  true,
				AttachStdout: true,
				AttachStderr: true,
				OpenStdin:    true,

				StdinOnce: true,
				Shell:     []string{"/bin/bash"},
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
	for key, value := range d.Args {
		if key == "command" {
			commandLine = append(commandLine, value.(string))
		}
		if key == "args" && value != nil {
			if args, ok := value.([]string); ok {
				for _, arg := range args {
					commandLine = append(commandLine, arg)
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

	output, _ := io.ReadAll(hijackedResponse.Conn)
	exitCode := execInspect.ExitCode

	if exitCode != 0 {
		toolCallFailureResponse.ToolCallResult.Error = errors.New(string(output))
		return toolCallFailureResponse
	}

	toolCallSuccessResponse.ToolCallResult.Output = string(output)

	return toolCallSuccessResponse
}

// TODO: implement the cleanup func.
func (d DockerExec) CleanUp() {
	/*
		STEPS:
		1. Stop all the running containers respectively.
		2. close all the connections to the docker client.
	*/
}
