package tools

import (
	"testing"
)

func TestExecute(t *testing.T) {

	dockerExec := DockerExec{
		Id:   "test-docker-exec",
		Name: DOCKER_EXEC,
		Args: map[string]any{
			"command": "nmap",
			"args":    []string{"--help"},
		},
	}

	msg := dockerExec.Execute("test-docker-exec")
	if msg.Status == Failure {
		t.Fatalf("ToolResponse status: %s\n", msg.ToolCallResult.Error)
	}
	t.Logf("output: %s", msg.ToolCallResult.Output)
}
