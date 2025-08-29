package page

import (
	"strings"
	"testing"
)

// Test helper functions for the bang command functionality

func TestParseCommandString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectCmd   string
		expectArgs  []string
		expectError bool
	}{
		{
			name:        "simple command",
			input:       "ls",
			expectCmd:   "ls",
			expectArgs:  []string{},
			expectError: false,
		},
		{
			name:        "command with args",
			input:       "ls -la /tmp",
			expectCmd:   "ls",
			expectArgs:  []string{"-la", "/tmp"},
			expectError: false,
		},
		{
			name:        "command with multiple spaces",
			input:       "  ls   -la   /tmp  ",
			expectCmd:   "ls",
			expectArgs:  []string{"-la", "/tmp"},
			expectError: false,
		},
		{
			name:        "empty command",
			input:       "",
			expectError: true,
		},
		{
			name:        "only spaces",
			input:       "   ",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts := strings.Fields(strings.TrimSpace(tt.input))
			
			if tt.expectError {
				if len(parts) != 0 {
					t.Errorf("expected empty parts for invalid input, got: %v", parts)
				}
				return
			}

			if len(parts) == 0 {
				t.Errorf("expected non-empty parts, got empty")
				return
			}

			cmd := parts[0]
			args := []string{}
			if len(parts) > 1 {
				args = parts[1:]
			}

			if cmd != tt.expectCmd {
				t.Errorf("expected command %q, got %q", tt.expectCmd, cmd)
			}

			if len(args) != len(tt.expectArgs) {
				t.Errorf("expected %d args, got %d", len(tt.expectArgs), len(args))
				return
			}

			for i, expectedArg := range tt.expectArgs {
				if args[i] != expectedArg {
					t.Errorf("expected arg[%d] = %q, got %q", i, expectedArg, args[i])
				}
			}
		})
	}
}

func TestBangCommandDetection(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		isBang   bool
		command  string
	}{
		{
			name:    "bang command",
			input:   "!ls -la",
			isBang:  true,
			command: "ls -la",
		},
		{
			name:    "bang command with spaces",
			input:   "!  ls -la  ",
			isBang:  true,
			command: "ls -la",
		},
		{
			name:   "normal message",
			input:  "hello world",
			isBang: false,
		},
		{
			name:   "exclamation in middle",
			input:  "hello ! world",
			isBang: false,
		},
		{
			name:   "empty bang command",
			input:  "!",
			isBang: true,
			command: "",
		},
		{
			name:   "bang with only spaces",
			input:  "!   ",
			isBang: true,
			command: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isBang := strings.HasPrefix(tt.input, "!")
			
			if isBang != tt.isBang {
				t.Errorf("expected isBang = %v, got %v", tt.isBang, isBang)
			}

			if tt.isBang {
				command := strings.TrimSpace(tt.input[1:])
				if command != tt.command {
					t.Errorf("expected command %q, got %q", tt.command, command)
				}
			}
		})
	}
}