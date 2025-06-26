package tools

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	docker "github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/yyovil/tandem/internal/utils"
)

type DockerExec Tool

var dockerClient *docker.Client
var kaliImage = "kali:withtools"
var shell types.HijackedResponse

var DockerExecTool = DockerExec{
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
	/*
		STEPS:
		2. spot a container if available, start one if not available assuming we already have the kali image.
		3. exec the command and pass the arguments to it.
		4. create a ToolResponse out of it and return it.
	*/
	var toolResponse ToolResponse
	if dockerClient == nil {
		var err error
		dockerClient, err = docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
		if err != nil {
			toolResponse = ToolResponse{
				Name:       DOCKER_EXEC,
				ToolCallId: toolCallId,
				Status:     Failure,
				Result: ToolResponseResult{
					Error: errors.Wrap(err, "failed to create docker client"),
				},
			}
			return toolResponse
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
		toolResponse = ToolResponse{
			Name:       DOCKER_EXEC,
			ToolCallId: toolCallId,
			Status:     Failure,
			Result: ToolResponseResult{
				Error: errors.Wrap(err, "failed to list docker containers"),
			},
		}
		return toolResponse
	}

	for _, summary := range summarySlice {
		// TODO: warn if bind mounts aren't set up. this maybe because the container is not created and ran by tandem.
		for _, mountPoint := range summary.Mounts {
			// NOTE: this works because we aren't going to have any other kind of mounts.
			if mountPoint.Type != mount.TypeBind {
				toolResponse = ToolResponse{
					Name:       DOCKER_EXEC,
					ToolCallId: toolCallId,
					Status:     Failure,
					Result: ToolResponseResult{
						Error: errors.New("bind mount not set up for container " + summary.ID),
					},
				}
				return toolResponse
			}
		}

		// after this switch expression, we assume that we have shell session for us.
		switch summary.Status {
		case "paused":
			// unpause it.
			if err := dockerClient.ContainerUnpause(ctx, summary.ID); err != nil {
				toolResponse = ToolResponse{
					Name:       DOCKER_EXEC,
					ToolCallId: toolCallId,
					Status:     Failure,
					Result: ToolResponseResult{
						Error: errors.Wrap(err, "failed to unpause docker container"),
					},
				}
				return toolResponse
			}

		case "exited",
			"created":
			{
				// TODO: start it.
				if err := dockerClient.ContainerStart(ctx, summary.ID, container.StartOptions{}); err != nil {
					toolResponse = ToolResponse{
						Name:       DOCKER_EXEC,
						ToolCallId: toolCallId,
						Status:     Failure,
						Result: ToolResponseResult{
							Error: errors.Wrap(err, "failed to start docker container"),
						},
					}
					return toolResponse
				}
			}
		default:
			// create a new container using the kali image.
			config := &container.Config{
				Image: kaliImage,
				Tty:   true,
				Cmd:   []string{"/bin/bash"},
				// continue... you were thinking about the stdin, stdout and stderr.
			}
			resp, err := dockerClient.ContainerCreate(ctx, config, nil, nil, nil, "")
			if err != nil {
				toolResponse = ToolResponse{
					Name:       DOCKER_EXEC,
					ToolCallId: toolCallId,
					Status:     Failure,
					Result: ToolResponseResult{
						Error: errors.Wrap(err, "failed to create docker container"),
					},
				}
				return toolResponse
			}

			if err := dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
				toolResponse = ToolResponse{
					Name:       DOCKER_EXEC,
					ToolCallId: toolCallId,
					Status:     Failure,
					Result: ToolResponseResult{
						Error: errors.Wrap(err, "failed to start docker container"),
					},
				}
				return toolResponse
			}
		}

		var commandLine []string
		for param, value := range d.Parameters {
			if param == "command" {
					commandLine = append(commandLine, value)
			}

			if param == "args" {

			}

		}
		shell.Conn.Write([]byte(strings.Join(commandLine, " ")))

	}

	return toolResponse
}
