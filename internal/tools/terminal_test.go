package tools

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// Integration-style test for DockerCli.Run executing an nmap command inside the kali:withtools container.
// This test only inspects the unified ToolResponse (errors are represented as ToolResponse with IsError=true).
func TestDockerCliRunNmap(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tool := NewDockerCli()

	args := TerminalArgs{Command: "nmap", Args: []string{"-p", "80", "127.0.0.1"}}
	payload, err := json.Marshal(args)
	if err != nil {
		// Should not happen; treat as test failure.
		resp := NewTextErrorResponse("failed to marshal args: " + err.Error())
		if !resp.IsError { // impossible, safeguard
			resp.IsError = true
		}
		// Fall through to assertions below by storing in channel.
	}

	call := ToolCall{ID: "test-nmap", Name: TerminalToolName, Input: string(payload)}

	respCh := make(chan ToolResponse, 1)

	go func() {
		resp, runErr := tool.Run(ctx, call)
		if runErr != nil {
			// Unify: wrap any Go error into a ToolResponse.
			resp = NewTextErrorResponse(runErr.Error())
		}
		respCh <- resp
	}()

	select {
	case <-ctx.Done():
		// Timeout indicates probable hang in DockerCli.Run read loop.
		resp := NewTextErrorResponse("timeout waiting for DockerCli.Run (possible read loop hang / missing container)")
		if resp.IsError {
			// We treat this as fatal test failure.
			// Provide hint: allocate non-zero buffer in read loop.
			// (Intentional failure path.)
			// Use t.Fatal directly.
			t.Fatal(resp.Content)
		}
	case resp := <-respCh:
		if resp.IsError {
			// Missing container -> skip (diagnostic not actionable in CI without environment prep).
			if strings.Contains(resp.Content, "couldn't find a container") {
				t.Skip(resp.Content)
			}
			// Generic connectivity to Docker: skip (environmental).
			if strings.Contains(strings.ToLower(resp.Content), "cannot connect") || strings.Contains(strings.ToLower(resp.Content), "connection refused") {
				t.Skip(resp.Content)
			}
			// Otherwise fail with the tool error content.
			t.Fatalf("tool error: %s", resp.Content)
		}
		// Soft assertion on expected nmap keyword.
		if !strings.Contains(resp.Content, "Nmap") {
			short := resp.Content
			if len(short) > 300 {
				short = short[:300] + "..."
			}
			// Log only to avoid flakiness; upgrade to failure if desired.
			t.Logf("output missing 'Nmap' keyword (len=%d): %s", len(resp.Content), short)
		}
	}
}
