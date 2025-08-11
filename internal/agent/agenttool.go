package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/yaydraco/tandem/internal/config"
	"github.com/yaydraco/tandem/internal/logging"
	"github.com/yaydraco/tandem/internal/message"
	"github.com/yaydraco/tandem/internal/session"
	"github.com/yaydraco/tandem/internal/subagent"
	"github.com/yaydraco/tandem/internal/tools"
)

const AgentToolName = "subagent"

var AgentNames = []string{
	string(config.Reconnoiter),
	string(config.VulnerabilityScanner),
	string(config.Exploiter),
	string(config.Reporter),
}

type AgentToolArgs struct {
	Prompt         string           `json:"prompt"`
	AgentName      config.AgentName `json:"agent_name,omitempty"`
	ExpectedOutput map[string]any   `json:"expected_output"`
}

type AgentTool struct {
	messages  message.Service
	sessions  session.Service
	subagents subagent.Service
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
				"description": "a JSON string representing the schema of the expected output that orchestrator requests the subagent to follow while responding after assigned task is completed.",
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

	// Validate agent name using slices.Contains
	if !slices.Contains(AgentNames, string(args.AgentName)) {
		return tools.NewTextErrorResponse("invalid agent name: " + string(args.AgentName)), nil
	}

	sessionID, messageID := tools.GetContextValues(ctx)
	if sessionID == "" || messageID == "" {
		return tools.ToolResponse{}, fmt.Errorf("session_id and message_id are required")
	}

	// NOTE: you can add more tools later here if needed on AgentName basis.
	agentTools := tools.PenetrationTestingAgentTools
	agent, err := NewAgent(args.AgentName, a.sessions, a.messages, agentTools, args.ExpectedOutput)
	if err != nil {
		return tools.NewTextErrorResponse("failed to create agent: " + err.Error()), nil
	}

	session, err := a.sessions.CreateTaskSession(ctx, call.ID, sessionID, fmt.Sprintf("%s agent's session", args.AgentName))
	if err != nil {
		return tools.ToolResponse{}, fmt.Errorf("error creating session: %s", err)
	}

	// Start tracking the activity
	activity, err := a.subagents.StartActivity(ctx, session.ID, sessionID, args.AgentName, args.Prompt)
	if err != nil {
		logging.Error("Failed to start activity tracking", err)
		// Continue without activity tracking if it fails
	}

	// Create a cancellable context for this specific task
	taskCtx, taskCancel := context.WithCancel(ctx)
	defer taskCancel()

	// Store the cancel function in the activity service
	if activity != nil {
		a.subagents.SetCancelFunc(activity.ID, taskCancel)
	}

	// Update activity status to running
	if activity != nil {
		a.subagents.UpdateActivity(ctx, activity.ID, subagent.StatusRunning, "Executing task...", "25%")
	}

	done, err := agent.Run(taskCtx, session.ID, args.Prompt)
	if err != nil {
		if activity != nil {
			a.subagents.CompleteActivity(ctx, activity.ID, false, fmt.Sprintf("error generating agent: %s", err))
		}
		return tools.ToolResponse{}, fmt.Errorf("error generating agent: %s", err)
	}

	logging.Debug("using agent", "name", args.AgentName, "busy", agent.IsBusy())
	
	// Update progress
	if activity != nil {
		a.subagents.UpdateActivity(ctx, activity.ID, subagent.StatusRunning, "Processing response...", "75%")
	}
	
	result := <-done
	logging.Debug("task done by agent", "name", args.AgentName, "busy", agent.IsBusy())
	
	if result.Error != nil {
		if activity != nil {
			if result.Error.Error() == "request cancelled by user" || result.Error.Error() == "context canceled" {
				// Activity was already marked as aborted by the cancel function
				return tools.ToolResponse{}, fmt.Errorf("task was aborted by user")
			}
			a.subagents.CompleteActivity(ctx, activity.ID, false, fmt.Sprintf("error generating agent: %s", result.Error))
		}
		return tools.ToolResponse{}, fmt.Errorf("error generating agent: %s", result.Error)
	}

	response := result.Message
	if response.Role != message.Assistant {
		if activity != nil {
			a.subagents.CompleteActivity(ctx, activity.ID, false, "no response")
		}
		return tools.NewTextErrorResponse("no response"), nil
	}

	updatedSession, err := a.sessions.Get(ctx, session.ID)
	if err != nil {
		if activity != nil {
			a.subagents.CompleteActivity(ctx, activity.ID, false, fmt.Sprintf("error getting session: %s", err))
		}
		return tools.ToolResponse{}, fmt.Errorf("error getting session: %s", err)
	}
	parentSession, err := a.sessions.Get(ctx, sessionID)
	if err != nil {
		if activity != nil {
			a.subagents.CompleteActivity(ctx, activity.ID, false, fmt.Sprintf("error getting parent session: %s", err))
		}
		return tools.ToolResponse{}, fmt.Errorf("error getting parent session: %s", err)
	}

	parentSession.Cost += updatedSession.Cost

	_, err = a.sessions.Save(ctx, parentSession)
	if err != nil {
		if activity != nil {
			a.subagents.CompleteActivity(ctx, activity.ID, false, fmt.Sprintf("error saving parent session: %s", err))
		}
		return tools.ToolResponse{}, fmt.Errorf("error saving parent session: %s", err)
	}
	
	// Mark activity as completed successfully
	if activity != nil {
		a.subagents.CompleteActivity(ctx, activity.ID, true, "Task completed successfully")
	}
	
	return tools.NewTextResponse(response.Content().String()), nil
}

func NewAgentTool(
	Sessions session.Service,
	Messages message.Service,
	SubAgents subagent.Service,
) tools.BaseTool {
	return &AgentTool{
		sessions:  Sessions,
		messages:  Messages,
		subagents: SubAgents,
	}
}
