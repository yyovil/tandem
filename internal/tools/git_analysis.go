package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	GitAnalysisToolName = "git_analysis"
)

type GitAnalysisParams struct {
	Repository string `json:"repository" description:"Path to the git repository (optional, defaults to current directory)"`
	CommitHash string `json:"commit_hash" description:"Specific commit hash to analyze (optional, defaults to latest commit)"`
	Branch     string `json:"branch" description:"Branch to analyze (optional, defaults to main/master)"`
	MaxCommits int    `json:"max_commits" description:"Maximum number of recent commits to analyze (optional, defaults to 5)"`
}

type GitAnalysisResult struct {
	Repository    string      `json:"repository"`
	Branch        string      `json:"branch"`
	LatestCommit  CommitInfo  `json:"latest_commit"`
	RecentCommits []CommitInfo `json:"recent_commits"`
	ChangedFiles  []FileChange `json:"changed_files"`
	Summary       string      `json:"summary"`
}

type CommitInfo struct {
	Hash      string    `json:"hash"`
	Author    string    `json:"author"`
	Date      time.Time `json:"date"`
	Message   string    `json:"message"`
	FilesChanged int    `json:"files_changed"`
}

type FileChange struct {
	Path      string `json:"path"`
	Status    string `json:"status"` // A=added, M=modified, D=deleted, R=renamed
	Language  string `json:"language"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
}

type GitAnalysisTool struct{}

func (g *GitAnalysisTool) Info() ToolInfo {
	return ToolInfo{
		Name:        GitAnalysisToolName,
		Description: "Analyze git repository commits and changes to understand recent development activity. Provides context for generating relevant demonstration content.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"repository": map[string]any{
					"type":        "string",
					"description": "Path to the git repository (optional, defaults to current directory)",
				},
				"commit_hash": map[string]any{
					"type":        "string",
					"description": "Specific commit hash to analyze (optional, defaults to latest commit)",
				},
				"branch": map[string]any{
					"type":        "string",
					"description": "Branch to analyze (optional, defaults to main/master)",
				},
				"max_commits": map[string]any{
					"type":        "integer",
					"description": "Maximum number of recent commits to analyze (optional, defaults to 5)",
				},
			},
			"required": []string{},
		},
		Required: []string{},
	}
}

func (g *GitAnalysisTool) Run(ctx context.Context, params ToolCall) (ToolResponse, error) {
	var gitParams GitAnalysisParams
	if err := json.Unmarshal([]byte(params.Input), &gitParams); err != nil {
		return NewTextErrorResponse(fmt.Sprintf("Invalid parameters: %v", err)), nil
	}

	// Set defaults
	if gitParams.Repository == "" {
		gitParams.Repository = "."
	}
	if gitParams.MaxCommits == 0 {
		gitParams.MaxCommits = 5
	}

	// Change to repository directory
	originalDir, err := os.Getwd()
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("Failed to get current directory: %v", err)), nil
	}

	if gitParams.Repository != "." {
		if err := os.Chdir(gitParams.Repository); err != nil {
			return NewTextErrorResponse(fmt.Sprintf("Failed to change to repository directory: %v", err)), nil
		}
		defer os.Chdir(originalDir)
	}

	// Check if we're in a git repository
	if !g.isGitRepository() {
		return NewTextErrorResponse("Not a git repository or no git repository found"), nil
	}

	// Get current branch if not specified
	if gitParams.Branch == "" {
		gitParams.Branch = g.getCurrentBranch()
	}

	// Analyze the repository
	result, err := g.analyzeRepository(ctx, gitParams)
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("Failed to analyze repository: %v", err)), nil
	}

	// Convert result to JSON for response
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return NewTextErrorResponse(fmt.Sprintf("Failed to marshal result: %v", err)), nil
	}

	return NewTextResponse(fmt.Sprintf("Git analysis completed successfully:\n%s", string(resultJSON))), nil
}

func (g *GitAnalysisTool) isGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

func (g *GitAnalysisTool) getCurrentBranch() string {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "main" // fallback
	}
	return strings.TrimSpace(string(output))
}

func (g *GitAnalysisTool) analyzeRepository(ctx context.Context, params GitAnalysisParams) (*GitAnalysisResult, error) {
	result := &GitAnalysisResult{
		Repository: params.Repository,
		Branch:     params.Branch,
	}

	// Get recent commits
	commits, err := g.getRecentCommits(ctx, params.MaxCommits)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent commits: %w", err)
	}
	result.RecentCommits = commits

	if len(commits) > 0 {
		result.LatestCommit = commits[0]

		// Get file changes for the latest commit
		changes, err := g.getFileChanges(ctx, commits[0].Hash)
		if err != nil {
			return nil, fmt.Errorf("failed to get file changes: %w", err)
		}
		result.ChangedFiles = changes

		// Generate summary
		result.Summary = g.generateSummary(commits, changes)
	}

	return result, nil
}

func (g *GitAnalysisTool) getRecentCommits(ctx context.Context, maxCommits int) ([]CommitInfo, error) {
	cmd := exec.CommandContext(ctx, "git", "log", 
		fmt.Sprintf("-%d", maxCommits),
		"--pretty=format:%H|%an|%ad|%s|%n",
		"--date=iso",
		"--numstat")
	
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git log command failed: %w", err)
	}

	return g.parseCommitLog(string(output))
}

func (g *GitAnalysisTool) parseCommitLog(output string) ([]CommitInfo, error) {
	var commits []CommitInfo
	lines := strings.Split(output, "\n")
	
	var currentCommit *CommitInfo
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Check if this is a commit info line (contains |)
		if strings.Contains(line, "|") && !strings.Contains(line, "\t") {
			parts := strings.Split(line, "|")
			if len(parts) >= 4 {
				// Save previous commit if exists
				if currentCommit != nil {
					commits = append(commits, *currentCommit)
				}
				
				// Parse date
				date, _ := time.Parse("2006-01-02 15:04:05 -0700", parts[2])
				
				currentCommit = &CommitInfo{
					Hash:    parts[0],
					Author:  parts[1],
					Date:    date,
					Message: parts[3],
				}
			}
		} else if currentCommit != nil && strings.Contains(line, "\t") {
			// This is a file change line
			currentCommit.FilesChanged++
		}
	}
	
	// Add the last commit
	if currentCommit != nil {
		commits = append(commits, *currentCommit)
	}
	
	return commits, nil
}

func (g *GitAnalysisTool) getFileChanges(ctx context.Context, commitHash string) ([]FileChange, error) {
	cmd := exec.CommandContext(ctx, "git", "show", "--numstat", "--name-status", commitHash)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git show command failed: %w", err)
	}

	return g.parseFileChanges(string(output))
}

func (g *GitAnalysisTool) parseFileChanges(output string) ([]FileChange, error) {
	var changes []FileChange
	lines := strings.Split(output, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Parse numstat format: additions	deletions	filename
		if strings.Contains(line, "\t") {
			parts := strings.Split(line, "\t")
			if len(parts) >= 3 {
				change := FileChange{
					Path:     parts[2],
					Language: g.detectLanguage(parts[2]),
				}
				
				// Parse additions/deletions
				if parts[0] != "-" {
					fmt.Sscanf(parts[0], "%d", &change.Additions)
				}
				if parts[1] != "-" {
					fmt.Sscanf(parts[1], "%d", &change.Deletions)
				}
				
				changes = append(changes, change)
			}
		}
	}
	
	return changes, nil
}

func (g *GitAnalysisTool) detectLanguage(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".go":
		return "go"
	case ".py":
		return "python"
	case ".js":
		return "javascript"
	case ".ts":
		return "typescript"
	case ".java":
		return "java"
	case ".c":
		return "c"
	case ".cpp", ".cxx", ".cc":
		return "cpp"
	case ".rs":
		return "rust"
	case ".sh":
		return "bash"
	case ".yml", ".yaml":
		return "yaml"
	case ".json":
		return "json"
	case ".md":
		return "markdown"
	case ".html":
		return "html"
	case ".css":
		return "css"
	case ".sql":
		return "sql"
	default:
		return "text"
	}
}

func (g *GitAnalysisTool) generateSummary(commits []CommitInfo, changes []FileChange) string {
	if len(commits) == 0 {
		return "No recent commits found"
	}

	latest := commits[0]
	summary := fmt.Sprintf("Latest commit: %s by %s\nMessage: %s\n", 
		latest.Hash[:8], latest.Author, latest.Message)

	if len(changes) > 0 {
		summary += fmt.Sprintf("Files changed: %d\n", len(changes))
		
		// Group by language
		langCount := make(map[string]int)
		for _, change := range changes {
			langCount[change.Language]++
		}
		
		summary += "Languages affected: "
		var langs []string
		for lang, count := range langCount {
			langs = append(langs, fmt.Sprintf("%s (%d files)", lang, count))
		}
		summary += strings.Join(langs, ", ")
	}

	return summary
}

// NewGitAnalysisTool creates a new Git analysis tool instance
func NewGitAnalysisTool() BaseTool {
	return &GitAnalysisTool{}
}