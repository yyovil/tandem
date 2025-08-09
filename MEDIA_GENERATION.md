# Automated Media Generation

Tandem now supports automated media generation through specialized AI agents that can analyze git commits and create demonstration content automatically.

## Overview

The media generation system consists of two specialized agents:

- **VHS Agent** (`vhs_agent`): Generates terminal demonstration videos using VHS (Video Hyper Scripts)
- **Freeze Agent** (`freeze_agent`): Creates beautiful SVG screenshots of code with syntax highlighting

## Features

### VHS Agent
- Analyzes git commits to understand new features and changes
- Generates VHS tape scripts that demonstrate functionality
- Creates engaging terminal recordings (GIFs) that showcase workflows
- Automatically configures appropriate terminal themes and settings

### Freeze Agent  
- Identifies important code changes from git commits
- Selects representative code snippets for demonstration
- Generates SVG images with syntax highlighting
- Ensures consistent styling and visual appeal

### Git Analysis Tool
- Analyzes repository commits and file changes
- Provides context about recent development activity
- Detects programming languages and categorizes changes
- Generates summaries for agent consumption

## Usage

### Manual Usage

You can interact with the media generation agents directly through Tandem:

```bash
# Generate VHS content for recent changes
./tandem --prompt "Analyze recent commits and create a VHS demo showing the new features"

# Generate Freeze screenshots for code changes  
./tandem --prompt "Create SVG screenshots of important code changes from the latest commit"
```

### Automated Usage (GitHub Actions)

The system automatically generates media content when changes are pushed to the main branch through the `.github/workflows/media-generation.yml` workflow.

#### Workflow Features:
- Triggers on pushes to main branch
- Can be manually triggered for specific commits
- Installs VHS and Freeze tools automatically
- Sets up virtual display for headless operation
- Generates media artifacts that can be downloaded
- Comments on PRs with generated content links

#### Skipping Media Generation
Add `[skip media]` to your commit message to skip automated media generation:

```bash
git commit -m "Fix typo [skip media]"
```

### Configuration

The agents are configured in `.tandem/swarm.json`:

```json
{
  "agents": {
    "vhs_agent": {
      "name": "VHS Media Generator",
      "model": "copilot.claude-sonnet-4",
      "tools": ["git_analysis", "vhs"]
    },
    "freeze_agent": {
      "name": "Freeze SVG Generator", 
      "model": "copilot.gpt-4o-mini",
      "tools": ["git_analysis", "freeze"]
    }
  }
}
```

## Tools

### VHS Tool
Generates terminal recordings using VHS tape scripts.

**Parameters:**
- `script`: VHS tape script content (required)
- `output_path`: Output GIF file path (default: "./demo.gif")
- `width`: Terminal width in pixels (default: 1200)
- `height`: Terminal height in pixels (default: 600)

### Freeze Tool
Creates SVG screenshots of code with syntax highlighting.

**Parameters:**
- `code`: Code content to capture (required)
- `language`: Programming language for syntax highlighting (required)
- `output_path`: Output SVG file path (default: "./code.svg")
- `theme`: Color theme (default: "catppuccin-frappe")
- `width`: Image width in pixels (default: 800)
- `font_family`: Font family (default: "JetBrains Mono")
- `font_size`: Font size in pixels (default: 14)

### Git Analysis Tool
Analyzes git repository commits and changes.

**Parameters:**
- `repository`: Repository path (default: current directory)
- `commit_hash`: Specific commit to analyze (optional)
- `branch`: Branch to analyze (default: main/master)
- `max_commits`: Maximum commits to analyze (default: 5)

## Requirements

### For Local Development
- Go 1.24+
- VHS tool: `https://github.com/charmbracelet/vhs`
- Freeze tool: `https://github.com/charmbracelet/freeze`
- API keys for AI providers (configured in `.env`)

### For GitHub Actions
- Repository secrets for API keys (e.g., `GROQ_API_KEY`)
- Ubuntu runner (tools are installed automatically)

## Installation

1. **Install VHS:**
   ```bash
   # macOS
   brew install vhs
   
   # Linux
   curl -s https://api.github.com/repos/charmbracelet/vhs/releases/latest | \
   grep "browser_download_url.*linux_amd64.tar.gz" | \
   cut -d '"' -f 4 | \
   xargs -I {} curl -L {} -o vhs.tar.gz && \
   tar -xzf vhs.tar.gz && \
   sudo mv vhs /usr/local/bin/
   ```

2. **Install Freeze:**
   ```bash
   # macOS
   brew install freeze
   
   # Linux  
   curl -s https://api.github.com/repos/charmbracelet/freeze/releases/latest | \
   grep "browser_download_url.*linux_amd64.tar.gz" | \
   cut -d '"' -f 4 | \
   xargs -I {} curl -L {} -o freeze.tar.gz && \
   tar -xzf freeze.tar.gz && \
   sudo mv freeze /usr/local/bin/
   ```

3. **Configure API Keys:**
   Copy `.env.example` to `.env` and add your API keys:
   ```bash
   cp .env.example .env
   # Edit .env to add your GROQ_API_KEY and other provider keys
   ```

## Examples

### VHS Script Generation
The VHS agent can generate scripts like:

```vhs
# Demo of new feature
Output feature_demo.gif
Set Width 1200
Set Height 600
Set Shell "bash"
Set Theme "Catppuccin Frappe"

Type "# Demonstrating new tandem feature"
Enter
Sleep 1s
Type "./tandem --help"
Enter
Sleep 2s
Type "# New functionality is now available!"
Enter
Sleep 1s
```

### Freeze Code Screenshots
The Freeze agent generates SVG images showing:
- New function implementations
- Configuration changes
- Important code additions
- API modifications

## Troubleshooting

### Common Issues

1. **VHS fails in headless environment:**
   - Ensure virtual display is set up (Xvfb)
   - Set `DISPLAY=:99` environment variable

2. **Freeze tool not found:**
   - Verify installation with `freeze --version`
   - Check PATH contains tool location

3. **No media generated:**
   - Check commit messages don't contain `[skip media]`
   - Verify API keys are configured
   - Check workflow logs for errors

### Debug Mode
Enable debug mode for detailed logging:

```bash
./tandem --debug --prompt "Your prompt here"
```

## Contributing

To add new media generation capabilities:

1. Create new tools in `internal/tools/`
2. Register tools in `tools.InitializeTools()`
3. Add agent configurations to `swarm.json`
4. Update JSON schema in `swarm.schema.json`
5. Test with the provided test script

## Security

- API keys are handled securely through environment variables
- Generated media files are temporary and cleaned up automatically
- No sensitive information is included in generated content