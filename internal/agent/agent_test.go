package agent

import (
	"context"
	"testing"

	"github.com/yyovil/tandem/internal/tools"
)

// MockTool is a simple mock tool for testing
type MockTool struct {
	name         string
	callCount    int
	shouldError  bool
	errorMessage string
}

func (m *MockTool) Info() tools.ToolInfo {
	return tools.ToolInfo{
		Name:        m.name,
		Description: "A mock tool for testing",
		Parameters:  map[string]any{},
		Required:    []string{},
	}
}

func (m *MockTool) Run(ctx context.Context, params tools.ToolCall) (tools.ToolResponse, error) {
	m.callCount++
	if m.shouldError {
		return tools.NewTextErrorResponse(m.errorMessage), nil
	}
	return tools.NewTextResponse("Mock tool executed successfully"), nil
}

func TestToolCallLimitEnforcement(t *testing.T) {
	tracker := tools.NewToolCallTracker()
	
	// Test that limits are respected
	sessionID := "test-session"
	toolName := "test-tool"
	maxCalls := 3

	// Test normal operation under limit
	for i := 0; i < maxCalls; i++ {
		count := tracker.GetToolCallCount(sessionID, toolName)
		if count >= maxCalls {
			t.Errorf("Tool call should be allowed. Current count: %d, max: %d", count, maxCalls)
		}
		tracker.IncrementToolCall(sessionID, toolName)
	}

	// Test limit exceeded
	count := tracker.GetToolCallCount(sessionID, toolName)
	if count < maxCalls {
		t.Errorf("Expected to reach limit. Current count: %d, max: %d", count, maxCalls)
	}

	// Verify final count
	if count != maxCalls {
		t.Errorf("Expected count to be %d, got %d", maxCalls, count)
	}
}

func TestToolCallTracker_MultipleTools(t *testing.T) {
	tracker := tools.NewToolCallTracker()
	sessionID := "test-session"
	
	// Test that different tools have independent counters
	tool1 := "tool1"
	tool2 := "tool2"
	
	tracker.IncrementToolCall(sessionID, tool1)
	tracker.IncrementToolCall(sessionID, tool1)
	tracker.IncrementToolCall(sessionID, tool2)
	
	if count := tracker.GetToolCallCount(sessionID, tool1); count != 2 {
		t.Errorf("Expected tool1 count to be 2, got %d", count)
	}
	
	if count := tracker.GetToolCallCount(sessionID, tool2); count != 1 {
		t.Errorf("Expected tool2 count to be 1, got %d", count)
	}
}

func TestToolCallTracker_MultipleSessions(t *testing.T) {
	tracker := tools.NewToolCallTracker()
	toolName := "test-tool"
	
	// Test that different sessions have independent counters
	session1 := "session1"
	session2 := "session2"
	
	tracker.IncrementToolCall(session1, toolName)
	tracker.IncrementToolCall(session1, toolName)
	tracker.IncrementToolCall(session2, toolName)
	
	if count := tracker.GetToolCallCount(session1, toolName); count != 2 {
		t.Errorf("Expected session1 count to be 2, got %d", count)
	}
	
	if count := tracker.GetToolCallCount(session2, toolName); count != 1 {
		t.Errorf("Expected session2 count to be 1, got %d", count)
	}
}