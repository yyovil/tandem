package tools

//TODO: make it compatible with genai.Tool
type Tool string

const (
	DockerContainerAvailable Tool = "docker_container_available"
	DockerContainerStart     Tool = "docker_container_start"
	DockerExec               Tool = "docker_exec"
)
