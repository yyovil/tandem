package page

import (
	"context"
	"fmt"
	"testing"

	"github.com/yyovil/tandem/internal/tools"
	"github.com/google/uuid"
)

// TestTerminalToolIntegration tests the terminal tool execution without requiring Docker
func TestTerminalToolIntegration(t *testing.T) {
	// This test verifies our terminal tool usage pattern works correctly
	// even if Docker container isn't available (it will fail gracefully)
	
	terminalTool := tools.NewDockerCli()
	
	// Test valid tool call structure with JSON directly
	argsJSON := `{"command":"echo","args":["hello","world"]}`
	
	toolCall := tools.ToolCall{
		ID:    uuid.New().String(),
		Name:  "terminal",
		Input: argsJSON,
	}
	
	// Execute the command - this will fail without Docker but should handle gracefully
	response, err := terminalTool.Run(context.Background(), toolCall)
	
	// We expect either a successful response or a graceful error
	if err != nil {
		t.Logf("Expected error without Docker container: %v", err)
	} else if response.IsError {
		t.Logf("Expected error response without Docker container: %s", response.Content)
	} else {
		t.Logf("Successful execution (Docker container available): %s", response.Content)
	}
	
	// Test the structure of our response
	if response.Type == "" {
		t.Error("ToolResponse should have a Type field set")
	}
}

// TestExecuteCommandErrorHandling tests error handling in our executeCommand function
func TestExecuteCommandErrorHandling(t *testing.T) {
	// Create a mock chatPage - we'll only test the executeCommand logic
	page := &chatPage{}
	
	// Test empty command
	_, err := page.executeCommand("")
	if err == nil {
		t.Error("Expected error for empty command")
	}
	
	// Test whitespace only command  
	_, err = page.executeCommand("   ")
	if err == nil {
		t.Error("Expected error for whitespace-only command")
	}
	
	// Test normal command structure (will fail without Docker but should parse correctly)
	_, err = page.executeCommand("ls -la")
	if err != nil {
		// This is expected without Docker container, but error should be about Docker, not parsing
		if fmt.Sprintf("%v", err) == "empty command" {
			t.Error("Command parsing failed - should not be 'empty command' error")
		}
		t.Logf("Expected Docker-related error: %v", err)
	}
}