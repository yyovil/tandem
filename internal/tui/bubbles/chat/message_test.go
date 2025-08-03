package chat

import (
	"strings"
	"testing"
	"time"

	"github.com/yaydraco/tandem/internal/message"
)

func TestFormatTimestampDiffSeconds(t *testing.T) {
	tests := []struct {
		name     string
		start    int64
		end      int64
		expected string
	}{
		{
			name:     "less than 1 second",
			start:    1000,
			end:      1500,
			expected: "0s",
		},
		{
			name:     "exactly 1 second",
			start:    1000,
			end:      2000,
			expected: "1s",
		},
		{
			name:     "multiple seconds",
			start:    1000,
			end:      3500,
			expected: "2s",
		},
		{
			name:     "exactly 60 seconds",
			start:    1000,
			end:      61000,
			expected: "1m",
		},
		{
			name:     "1 minute 23 seconds",
			start:    1000,
			end:      84000,
			expected: "1m23s",
		},
		{
			name:     "exactly 1 hour",
			start:    1000,
			end:      3601000,
			expected: "1h",
		},
		{
			name:     "1 hour 20 minutes",
			start:    1000,
			end:      4801000,
			expected: "1h20m",
		},
		{
			name:     "34 seconds",
			start:    1000,
			end:      35000,
			expected: "34s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTimestampDiffSeconds(tt.start, tt.end)
			if result != tt.expected {
				t.Errorf("formatTimestampDiffSeconds(%d, %d) = %s, want %s", tt.start, tt.end, result, tt.expected)
			}
		})
	}
}

func TestFormatTimeFromTimestamp(t *testing.T) {
	// Test with a known timestamp (January 1, 2024 3:45 PM UTC)
	timestamp := int64(1704118500000) // milliseconds
	result := formatTimeFromTimestamp(timestamp)

	// Parse it to ensure it's a valid time
	_, err := time.Parse("3:04 PM", result)
	if err != nil {
		t.Errorf("formatTimeFromTimestamp() returned invalid time format: %s, error: %v", result, err)
	}

	// Check that it ends with AM or PM
	if !strings.HasSuffix(result, " AM") && !strings.HasSuffix(result, " PM") {
		t.Errorf("formatTimeFromTimestamp() = %s, want format ending with AM or PM", result)
	}

	// Check that the result has the correct format (H:MM AM/PM or HH:MM AM/PM)
	if len(result) < 7 || len(result) > 8 {
		t.Errorf("formatTimeFromTimestamp() length = %d, want 7 or 8", len(result))
	}
}

func TestRenderAssistantMessage_TimeDisplay(t *testing.T) {
	// Create a test message with finish data using fixed timestamps
	startTime := int64(1704118500000) // January 1, 2024 3:45 PM UTC in milliseconds
	endTime := startTime + 2500       // 2.5 seconds later

	msg := message.Message{
		ID:        "test-id",
		CreatedAt: startTime,
		Parts: []message.ContentPart{
			message.TextContent{Text: "Hello, this is a test message"},
			message.Finish{
				Reason: message.FinishReasonEndTurn,
				Time:   endTime,
			},
		},
	}

	// Render the message
	uiMessages := renderAssistantMessage(msg, []message.Message{}, false, 80, 0)

	// There should be one UI message
	if len(uiMessages) != 1 {
		t.Errorf("Expected 1 UI message, got %d", len(uiMessages))
		return
	}

	content := uiMessages[0].content

	// Check that the content contains both timestamp format and duration
	// The format should be "H:MM AM/PM (Xs)" or "HH:MM AM/PM (Xs)"
	if !strings.Contains(content, ":") {
		t.Errorf("Message content should contain time with colons, got: %s", content)
	}

	if !strings.Contains(content, "M") {
		t.Errorf("Message content should contain AM or PM, got: %s", content)
	}

	if !strings.Contains(content, "s") {
		t.Errorf("Message content should contain 's' for seconds, got: %s", content)
	}

	if !strings.Contains(content, "(") || !strings.Contains(content, ")") {
		t.Errorf("Message content should contain duration in parentheses, got: %s", content)
	}
}

func TestRenderAssistantMessage_NonFinishedMessage(t *testing.T) {
	// Create a test message without finish data
	msg := message.Message{
		ID:        "test-id",
		CreatedAt: time.Now().Unix() * 1000, // Convert to milliseconds
		Parts: []message.ContentPart{
			message.TextContent{Text: "Hello, this is an unfinished message"},
		},
	}

	// Render the message
	uiMessages := renderAssistantMessage(msg, []message.Message{}, false, 80, 0)

	// There should be one UI message
	if len(uiMessages) != 1 {
		t.Errorf("Expected 1 UI message, got %d", len(uiMessages))
		return
	}

	content := uiMessages[0].content

	// Check that the content does not contain time format for unfinished messages
	// Should not contain ":" from time format or "(" from duration
	lines := strings.Split(content, "\n")
	lastLine := lines[len(lines)-1]

	// The time format should not appear in an unfinished message
	if strings.Contains(lastLine, ":") && strings.Contains(lastLine, "(") {
		t.Errorf("Unfinished message should not contain time display, got: %s", lastLine)
	}
}
