package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	FreezeToolName = "freeze"
)

type FreezeParams struct {
	Code       string `json:"code" description:"The code content to capture as SVG"`
	Language   string `json:"language" description:"Programming language for syntax highlighting (e.g., 'go', 'bash', 'python')"`
	OutputPath string `json:"output_path" description:"Path where the generated SVG should be saved (optional, defaults to ./code.svg)"`
	Theme      string `json:"theme" description:"Color theme to use (optional, defaults to 'catppuccin-frappe')"`
	Width      int    `json:"width" description:"Width of the output in pixels (optional, defaults to 800)"`
	Height     int    `json:"height" description:"Height of the output in pixels (optional, auto-calculated if not specified)"`
	FontFamily string `json:"font_family" description:"Font family to use (optional, defaults to 'JetBrains Mono')"`
	FontSize   int    `json:"font_size" description:"Font size in pixels (optional, defaults to 14)"`
	ShowLineNumbers bool `json:"show_line_numbers" description:"Whether to show line numbers (optional, defaults to true)"`
	WindowFrame bool `json:"window_frame" description:"Whether to show window frame decoration (optional, defaults to true)"`
}

type FreezeTool struct{}

func (f *FreezeTool) Info() ToolInfo {
	return ToolInfo{
		Name:        FreezeToolName,
		Description: "Generate beautiful SVG screenshots of code using Freeze. Creates syntax-highlighted, styled code images perfect for documentation and demos.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"code": map[string]any{
					"type":        "string",
					"description": "The code content to capture as SVG",
				},
				"language": map[string]any{
					"type":        "string",
					"description": "Programming language for syntax highlighting (e.g., 'go', 'bash', 'python')",
				},
				"output_path": map[string]any{
					"type":        "string",
					"description": "Path where the generated SVG should be saved (optional, defaults to ./code.svg)",
				},
				"theme": map[string]any{
					"type":        "string",
					"description": "Color theme to use (optional, defaults to 'catppuccin-frappe')",
				},
				"width": map[string]any{
					"type":        "integer",
					"description": "Width of the output in pixels (optional, defaults to 800)",
				},
				"height": map[string]any{
					"type":        "integer",
					"description": "Height of the output in pixels (optional, auto-calculated if not specified)",
				},
				"font_family": map[string]any{
					"type":        "string",
					"description": "Font family to use (optional, defaults to 'JetBrains Mono')",
				},
				"font_size": map[string]any{
					"type":        "integer",
					"description": "Font size in pixels (optional, defaults to 14)",
				},
				"show_line_numbers": map[string]any{
					"type":        "boolean",
					"description": "Whether to show line numbers (optional, defaults to true)",
				},
				"window_frame": map[string]any{
					"type":        "boolean",
					"description": "Whether to show window frame decoration (optional, defaults to true)",
				},
			},
			"required": []string{"code", "language"},
		},
		Required: []string{"code", "language"},
	}
}

func (f *FreezeTool) Run(ctx context.Context, params ToolCall) (ToolResponse, error) {
	var freezeParams FreezeParams
	if err := json.Unmarshal([]byte(params.Input), &freezeParams); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("Invalid parameters: %v", err)), nil
	}

	if freezeParams.Code == "" {
		return NewTextErrorResponse("Code content is required"), nil
	}
	if freezeParams.Language == "" {
		return NewTextErrorResponse("Language is required"), nil
	}

	// Set defaults
	if freezeParams.OutputPath == "" {
		freezeParams.OutputPath = "./code.svg"
	}
	if freezeParams.Theme == "" {
		freezeParams.Theme = "catppuccin-frappe"
	}
	if freezeParams.Width == 0 {
		freezeParams.Width = 800
	}
	if freezeParams.FontFamily == "" {
		freezeParams.FontFamily = "JetBrains Mono"
	}
	if freezeParams.FontSize == 0 {
		freezeParams.FontSize = 14
	}

	// Check if Freeze is installed
	if err := checkFreezeInstalled(); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("Freeze is not installed: %v", err)), nil
	}

	// Create temporary code file
	codeFile, err := f.createCodeFile(freezeParams)
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("Failed to create code file: %v", err)), nil
	}
	defer os.Remove(codeFile)

	// Execute Freeze
	output, err := f.executeFreeze(ctx, codeFile, freezeParams)
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("Failed to execute Freeze: %v\nOutput: %s", err, output)), nil
	}

	// Verify output file was created
	if _, err := os.Stat(freezeParams.OutputPath); os.IsNotExist(err) {
		return NewTextErrorResponse(fmt.Sprintf("Output file was not created: %s", freezeParams.OutputPath)), nil
	}

	return NewTextResponse(fmt.Sprintf("Freeze SVG successfully generated: %s\nCode captured:\n%s", freezeParams.OutputPath, truncateString(freezeParams.Code, 200))), nil
}

func (f *FreezeTool) createCodeFile(params FreezeParams) (string, error) {
	// Create a temporary file with appropriate extension
	timestamp := time.Now().Format("20060102_150405")
	ext := getFileExtensionForLanguage(params.Language)
	codeFile := fmt.Sprintf("/tmp/freeze_code_%s.%s", timestamp, ext)

	if err := os.WriteFile(codeFile, []byte(params.Code), 0644); err != nil {
		return "", fmt.Errorf("failed to write code file: %w", err)
	}

	return codeFile, nil
}

func (f *FreezeTool) executeFreeze(ctx context.Context, codeFile string, params FreezeParams) (string, error) {
	args := []string{
		codeFile,
		"--output", params.OutputPath,
		"--theme", params.Theme,
		"--width", fmt.Sprintf("%d", params.Width),
		"--font-family", params.FontFamily,
		"--font-size", fmt.Sprintf("%d", params.FontSize),
	}

	if params.Height > 0 {
		args = append(args, "--height", fmt.Sprintf("%d", params.Height))
	}

	if params.ShowLineNumbers {
		args = append(args, "--line-numbers")
	}

	if params.WindowFrame {
		args = append(args, "--window")
	}

	cmd := exec.CommandContext(ctx, "freeze", args...)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("freeze command failed: %w", err)
	}

	return string(output), nil
}

func checkFreezeInstalled() error {
	cmd := exec.Command("freeze", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("freeze command not found. Please install Freeze from https://github.com/charmbracelet/freeze")
	}
	return nil
}

func getFileExtensionForLanguage(language string) string {
	langMap := map[string]string{
		"go":         "go",
		"python":     "py",
		"javascript": "js",
		"typescript": "ts",
		"java":       "java",
		"c":          "c",
		"cpp":        "cpp",
		"rust":       "rs",
		"bash":       "sh",
		"shell":      "sh",
		"yaml":       "yml",
		"json":       "json",
		"xml":        "xml",
		"html":       "html",
		"css":        "css",
		"sql":        "sql",
		"dockerfile": "dockerfile",
		"markdown":   "md",
	}

	if ext, exists := langMap[strings.ToLower(language)]; exists {
		return ext
	}
	return "txt"
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// NewFreezeTool creates a new Freeze tool instance
func NewFreezeTool() BaseTool {
	return &FreezeTool{}
}