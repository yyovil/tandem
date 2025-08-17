package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/yyovil/tandem/internal/models"
)

func TestGetAgentPrompt_Orchestrator(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "tandem_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a temporary .tandem directory
	tandemDir := filepath.Join(tempDir, ".tandem")
	err = os.MkdirAll(tandemDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .tandem dir: %v", err)
	}

	// Test cases
	testCases := []struct {
		name                string
		setupRoEFile        bool
		roEContent          string
		expectedInPrompt    []string
		notExpectedInPrompt []string
	}{
		{
			name:         "orchestrator without RoE file",
			setupRoEFile: false,
			expectedInPrompt: []string{
				"<description>",
				"<goal>",
				"<instructions>",
			},
			notExpectedInPrompt: []string{
				"<context>",
				"Rules of Engagement",
			},
		},
		{
			name:         "orchestrator with RoE file",
			setupRoEFile: true,
			roEContent:   "# Rules of Engagement\n\nThis is a test penetration test.\n\n## Scope\n- Target: test.example.com\n- Authorization: Written approval received",
			expectedInPrompt: []string{
				"<description>",
				"<goal>",
				"<instructions>",
				"<context>",
				"Rules of Engagement",
				"test.example.com",
				"Written approval received",
			},
			notExpectedInPrompt: []string{},
		},
		{
			name:         "orchestrator with empty RoE file",
			setupRoEFile: true,
			roEContent:   "",
			expectedInPrompt: []string{
				"<description>",
				"<goal>",
				"<instructions>",
				"<context>",
				"</context>",
			},
			notExpectedInPrompt: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset the global config and context
			cfg = nil
			contextContent = ""
			onceContext = sync.Once{}

			// Setup RoE file if needed
			roEPath := filepath.Join(tandemDir, "RoE.md")
			if tc.setupRoEFile {
				err := os.WriteFile(roEPath, []byte(tc.roEContent), 0644)
				if err != nil {
					t.Fatalf("Failed to write RoE file: %v", err)
				}
			}

			// Create a test configuration
			cfg = &Config{
				RoEPath:    roEPath,
				WorkingDir: tempDir,
				Agents: map[AgentName]Agent{
					Orchestrator: {
						AgentID:     "orchestrator",
						Name:        "Test Orchestrator",
						Description: "Test orchestrator agent for penetration testing coordination",
						Goal:        "Coordinate and manage penetration testing activities",
						Instructions: []string{
							"Coordinate with other agents",
							"Manage the overall testing workflow",
							"Ensure all testing follows the rules of engagement",
						},
					},
				},
			}

			// Test GetAgentPrompt for orchestrator
			prompt := GetAgentPrompt(Orchestrator, models.ProviderOpenAI)

			// Verify the prompt is not empty
			if prompt == "" {
				t.Error("Expected non-empty prompt")
			}

			// Check for expected content
			for _, expected := range tc.expectedInPrompt {
				if !strings.Contains(prompt, expected) {
					t.Errorf("Expected prompt to contain %q, but it didn't.\nPrompt: %s", expected, prompt)
				}
			}

			// Check for content that shouldn't be there
			for _, notExpected := range tc.notExpectedInPrompt {
				if strings.Contains(prompt, notExpected) {
					t.Errorf("Expected prompt to NOT contain %q, but it did.\nPrompt: %s", notExpected, prompt)
				}
			}

			// Verify the basic structure is correct
			if !strings.Contains(prompt, "Test orchestrator agent for penetration testing coordination") {
				t.Error("Expected prompt to contain agent description")
			}

			if !strings.Contains(prompt, "Coordinate and manage penetration testing activities") {
				t.Error("Expected prompt to contain agent goal")
			}

			if !strings.Contains(prompt, "Coordinate with other agents") {
				t.Error("Expected prompt to contain agent instructions")
			}
		})
	}
}

func TestGetAgentPrompt_Orchestrator_EdgeCases(t *testing.T) {
	// Test with nil config
	cfg = nil
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when config is nil")
		}
	}()
	GetAgentPrompt(Orchestrator, models.ProviderOpenAI)
}

func TestGetAgentPrompt_Orchestrator_DifferentProviders(t *testing.T) {
	// Reset the global config
	cfg = &Config{
		RoEPath:    ".tandem/RoE.md",
		WorkingDir: "/tmp",
		Agents: map[AgentName]Agent{
			Orchestrator: {
				AgentID:      "orchestrator",
				Name:         "Test Orchestrator",
				Description:  "Test orchestrator agent",
				Goal:         "Test goal",
				Instructions: []string{"Test instruction"},
			},
		},
	}

	// Test with different providers
	providers := []models.ModelProvider{
		models.ProviderOpenAI,
		models.ProviderAnthropic,
		models.ProviderGemini,
		models.ProviderCopilot,
	}

	for _, provider := range providers {
		t.Run(fmt.Sprintf("provider_%s", provider), func(t *testing.T) {
			prompt := GetAgentPrompt(Orchestrator, provider)

			if prompt == "" {
				t.Errorf("Expected non-empty prompt for provider %s", provider)
			}

			// The prompt should contain the same basic structure regardless of provider
			if !strings.Contains(prompt, "Test orchestrator agent") {
				t.Errorf("Expected prompt to contain description for provider %s", provider)
			}

			if !strings.Contains(prompt, "Test goal") {
				t.Errorf("Expected prompt to contain goal for provider %s", provider)
			}
		})
	}
}

func TestGetAgentPrompt_Orchestrator_RoEFileHandling(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "tandem_roe_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test with non-existent RoE file
	t.Run("non_existent_roe_file", func(t *testing.T) {
		cfg = nil
		contextContent = ""
		onceContext = sync.Once{}

		nonExistentPath := filepath.Join(tempDir, "non_existent.md")
		cfg = &Config{
			RoEPath:    nonExistentPath,
			WorkingDir: tempDir,
			Agents: map[AgentName]Agent{
				Orchestrator: {
					AgentID:      "orchestrator",
					Name:         "Test Orchestrator",
					Description:  "Test description",
					Goal:         "Test goal",
					Instructions: []string{"Test instruction"},
				},
			},
		}

		prompt := GetAgentPrompt(Orchestrator, models.ProviderOpenAI)

		// Should not contain context section when file doesn't exist
		if strings.Contains(prompt, "<context>") {
			t.Error("Expected no context section when RoE file doesn't exist")
		}
	})

	// Test with RoE file that has complex content
	t.Run("complex_roe_content", func(t *testing.T) {
		cfg = nil
		contextContent = ""
		onceContext = sync.Once{}

		roEPath := filepath.Join(tempDir, "complex_RoE.md")
		complexContent := `# Penetration Testing Rules of Engagement

## Scope
- **Target Systems**: web.example.com, api.example.com
- **IP Ranges**: 192.168.1.0/24, 10.0.0.0/8
- **Excluded Systems**: production.example.com

## Authorization
- Written authorization received from CISO
- Testing window: 2025-01-15 to 2025-01-20

## Constraints
- No DoS attacks
- No data exfiltration
- Business hours only: 9 AM - 5 PM EST

## Reporting
- Daily status updates required
- Final report due within 48 hours of completion
`

		err := os.WriteFile(roEPath, []byte(complexContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write complex RoE file: %v", err)
		}

		cfg = &Config{
			RoEPath:    roEPath,
			WorkingDir: tempDir,
			Agents: map[AgentName]Agent{
				Orchestrator: {
					AgentID:      "orchestrator",
					Name:         "Test Orchestrator",
					Description:  "Test description",
					Goal:         "Test goal",
					Instructions: []string{"Test instruction"},
				},
			},
		}

		prompt := GetAgentPrompt(Orchestrator, models.ProviderOpenAI)

		// Should contain context section and specific content
		expectedContent := []string{
			"<context>",
			"web.example.com",
			"192.168.1.0/24",
			"Written authorization received",
			"No DoS attacks",
			"Daily status updates",
		}

		for _, expected := range expectedContent {
			if !strings.Contains(prompt, expected) {
				t.Errorf("Expected prompt to contain %q, but it didn't", expected)
			}
		}
	})
}
