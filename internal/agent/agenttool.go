package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/yaydraco/tandem/internal/config"
	"github.com/yaydraco/tandem/internal/message"
	"github.com/yaydraco/tandem/internal/session"
	"github.com/yaydraco/tandem/internal/tools"
)

const AgentToolName = "agent_tool"

var AgentNames = []config.AgentName{
	config.Reconnoiter,
	config.VulnerabilityScanner,
	config.Exploiter,
	config.Reporter,
}

type AgentToolArgs struct {
	Prompt         string             `json:"prompt"`
	AgentName      config.AgentName   `json:"agent_name,omitempty"`
	ExpectedOutput ExpectedOutputType `json:"expected_output,omitempty"`
}

type ExpectedOutputType struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Schema      map[string]any `json:"schema"`
}

type AgentTool struct {
	messages message.Service
	sessions session.Service
}

func (a *AgentTool) Info() tools.ToolInfo {
	return tools.ToolInfo{
		Name:        AgentToolName,
		Description: "A tool to assign penetration testing engagement related tasks to agents as per their area of expertise and purview.",
		Parameters: map[string]any{
			"prompt": map[string]any{
				"type":        "string",
				"description": "Prompt to send to the agent",
			},
			"agent_name": map[string]any{
				"type":        "string",
				"description": "ID of the agent to call",
				"enum":        AgentNames,
			},
			"expected_output": map[string]any{
				"type":        "object",
				"description": "Expected output schema specification that the agent should follow when generating structured responses.",
				// TODO: this has to adhere to this schema as well.
				"properties": map[string]any{
					"type": "object",
					"name": map[string]any{
						"type":        "string",
						"description": "Name identifier for the output schema",
					},
					"description": map[string]any{
						"type":        "string",
						"description": "Human-readable description of what the output should contain",
					},
					"schema": map[string]any{
						"type":        "object",
						"description": "JSON schema definition that specifies the structure, types, and constraints for the expected output",
					},
				},
				"required": []string{"name", "description", "schema"},
			},
		},
		Required: []string{"prompt", "agent_name", "expected_output"},
	}
}

func (a *AgentTool) Run(ctx context.Context, call tools.ToolCall) (tools.ToolResponse, error) {
	var args AgentToolArgs
	if err := json.Unmarshal([]byte(call.Input), &args); err != nil {
		return tools.NewTextErrorResponse("failed to parse agent tool parameters: " + err.Error()), nil
	}

	// Validate the agent name
	for _, agentName := range AgentNames {
		if args.AgentName == agentName {
			break
		}
		return tools.NewTextErrorResponse("invalid agent name: " + string(args.AgentName)), nil
	}

	sessionID, messageID := tools.GetContextValues(ctx)
	if sessionID == "" || messageID == "" {
		return tools.ToolResponse{}, fmt.Errorf("session_id and message_id are required")
	}

	// NOTE: you can add more tools later here if needed on AgentName basis.
	agentTools := tools.PenetrationTestingAgentTools
	agent, err := NewAgent(args.AgentName, a.sessions, a.messages, agentTools)
	if err != nil {
		return tools.NewTextErrorResponse("failed to create agent: " + err.Error()), nil
	}

	session, err := a.sessions.CreateTaskSession(ctx, call.ID, sessionID, "New Agent Session")
	if err != nil {
		return tools.ToolResponse{}, fmt.Errorf("error creating session: %s", err)
	}

	done, err := agent.Run(ctx, session.ID, args.Prompt)
	if err != nil {
		return tools.ToolResponse{}, fmt.Errorf("error generating agent: %s", err)
	}
	result := <-done
	if result.Error != nil {
		return tools.ToolResponse{}, fmt.Errorf("error generating agent: %s", result.Error)
	}

	response := result.Message
	if response.Role != message.Assistant {
		return tools.NewTextErrorResponse("no response"), nil
	}

	updatedSession, err := a.sessions.Get(ctx, session.ID)
	if err != nil {
		return tools.ToolResponse{}, fmt.Errorf("error getting session: %s", err)
	}
	parentSession, err := a.sessions.Get(ctx, sessionID)
	if err != nil {
		return tools.ToolResponse{}, fmt.Errorf("error getting parent session: %s", err)
	}

	parentSession.Cost += updatedSession.Cost

	_, err = a.sessions.Save(ctx, parentSession)
	if err != nil {
		return tools.ToolResponse{}, fmt.Errorf("error saving parent session: %s", err)
	}
	return tools.NewTextResponse(response.Content().String()), nil
}

/*
NOTE: we havent passed the expected output schema to the agent yet.
*/

func NewAgentTool(
	Sessions session.Service,
	Messages message.Service,
) tools.BaseTool {
	return &AgentTool{
		sessions: Sessions,
		messages: Messages,
	}
}
