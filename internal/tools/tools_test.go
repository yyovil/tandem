package tools

import (
	"testing"
)

func TestToolCallTracker(t *testing.T) {
	tracker := NewToolCallTracker()

	sessionID := "test-session"
	toolName := "test-tool"

	// Test initial count
	if count := tracker.GetToolCallCount(sessionID, toolName); count != 0 {
		t.Errorf("Expected initial count to be 0, got %d", count)
	}

	// Test incrementing
	tracker.IncrementToolCall(sessionID, toolName)
	if count := tracker.GetToolCallCount(sessionID, toolName); count != 1 {
		t.Errorf("Expected count to be 1 after increment, got %d", count)
	}

	// Test multiple increments
	tracker.IncrementToolCall(sessionID, toolName)
	tracker.IncrementToolCall(sessionID, toolName)
	if count := tracker.GetToolCallCount(sessionID, toolName); count != 3 {
		t.Errorf("Expected count to be 3 after multiple increments, got %d", count)
	}

	// Test different tool in same session
	anotherTool := "another-tool"
	tracker.IncrementToolCall(sessionID, anotherTool)
	if count := tracker.GetToolCallCount(sessionID, anotherTool); count != 1 {
		t.Errorf("Expected count for another tool to be 1, got %d", count)
	}

	// Verify first tool count unchanged
	if count := tracker.GetToolCallCount(sessionID, toolName); count != 3 {
		t.Errorf("Expected count for first tool to remain 3, got %d", count)
	}

	// Test different session
	anotherSession := "another-session"
	tracker.IncrementToolCall(anotherSession, toolName)
	if count := tracker.GetToolCallCount(anotherSession, toolName); count != 1 {
		t.Errorf("Expected count for tool in another session to be 1, got %d", count)
	}

	// Test reset session
	tracker.ResetSession(sessionID)
	if count := tracker.GetToolCallCount(sessionID, toolName); count != 0 {
		t.Errorf("Expected count to be 0 after session reset, got %d", count)
	}
	if count := tracker.GetToolCallCount(sessionID, anotherTool); count != 0 {
		t.Errorf("Expected count for another tool to be 0 after session reset, got %d", count)
	}

	// Verify other session unchanged
	if count := tracker.GetToolCallCount(anotherSession, toolName); count != 1 {
		t.Errorf("Expected count in another session to remain 1 after resetting different session, got %d", count)
	}
}

func TestToolCallTracker_EdgeCases(t *testing.T) {
	tracker := NewToolCallTracker()

	// Test with empty session ID
	emptySession := ""
	toolName := "test-tool"
	tracker.IncrementToolCall(emptySession, toolName)
	if count := tracker.GetToolCallCount(emptySession, toolName); count != 1 {
		t.Errorf("Expected count with empty session ID to be 1, got %d", count)
	}

	// Test with empty tool name
	sessionID := "test-session"
	emptyTool := ""
	tracker.IncrementToolCall(sessionID, emptyTool)
	if count := tracker.GetToolCallCount(sessionID, emptyTool); count != 1 {
		t.Errorf("Expected count with empty tool name to be 1, got %d", count)
	}

	// Test reset non-existent session
	tracker.ResetSession("non-existent-session")
	// Should not panic and other data should remain intact
	if count := tracker.GetToolCallCount(sessionID, emptyTool); count != 1 {
		t.Errorf("Expected existing data to remain after resetting non-existent session, got %d", count)
	}
}