package tools

import "github.com/yyovil/tandem/internal/utils"

type ToolName string

const (
	DOCKER_EXEC ToolName = "docker_exec"
)

type ToolParameters map[string]Param

// NOTE: use this type when defining tools for the agent.
type Tool struct {
	Description string
	Name        ToolName
	Parameters  ToolParameters
	Required    []string
}

type RunTool interface {
	Execute(toolCallId string) ToolResponse
}

// NOTE: use this type when defining tool calls from the agent response.
type ToolCall struct {
	Id   string
	Name ToolName
	Args map[string]any
}

type Status string

const (
	Success Status = "success"
	Failure Status = "failure"
)

// NOTE: use this type when defining tool results from the agent response.
type ToolResponse struct {
	Status     Status
	ToolCallId string
	Name       ToolName
	Result     ToolResponseResult
}

type ToolResponseResult struct {
	Error  error
	Output map[string]any
}

var ToolRegistry = map[ToolName]RunTool{
	DOCKER_EXEC: DockerExecTool,
}

/*
NOTE: trying to trick the go compiler here to get past the typing system. also this never gets executed because we never going to have a tool in the registry with the value of type Tool. every tool will have its own type alias and it has to have a type alias because they have to implement their own kind of Execute method.
*/
func (t Tool) Execute(toolCallId string) ToolResponse {
	return ToolResponse{}
}

func GetTool(name ToolName) (Tool, bool) {
	tool, ok := ToolRegistry[name]
	return tool.(Tool), ok
}

type Param struct {
	Description string
	Type        utils.Type     // TODO: make it more restrictive to OpenAPI 3.0 types
	Properties  ToolParameters // NOTE: properties are only required when used with Type "object".
	Items       *Param         // NOTE: items are only required when used with Type "array" any ways.
}
