# SubAgent Activity Monitoring

This feature allows users to view real-time activity of subagents and abort tasks when needed.

## Overview

The SubAgent Activity Monitoring system provides:

- **Real-time Activity Monitoring**: View what each subagent is doing in real-time
- **Intelligent Status Updates**: LLM-generated status messages that are context-aware
- **Task Abortion**: Cancel running subagent tasks with a hotkey
- **Progress Tracking**: Visual progress indicators with estimated time remaining
- **Multi-agent Support**: Track multiple concurrent subagent activities

## Usage

### Accessing the Activity Monitor

Press `Ctrl+T` from the main chat interface to open the SubAgent Activity page.

### Activity Page Features

The activity page displays a table with the following columns:

- **Agent**: The type of subagent (Reconnoiter, VulnerabilityScanner, Exploiter, Reporter)
- **Task**: Brief description of the assigned task
- **Status**: Current status with visual indicators and intelligent descriptions
- **Progress**: Percentage completion
- **ETA**: Estimated time remaining
- **Started**: Time when the task began
- **Duration**: How long the task has been running

### Visual Status Indicators

- 🔄 **Starting**: Task is initializing
- ⚡ **Running**: Task is actively executing
- ✅ **Completed**: Task finished successfully
- ❌ **Error**: Task encountered an error
- 🛑 **Aborted**: Task was cancelled by user

### Keyboard Controls

- `Ctrl+T`: Open SubAgent Activity page
- `Ctrl+A`: Abort selected task
- `R`: Refresh activity list
- `Esc`: Return to chat page

## Intelligent Status Messages

The system generates context-aware status messages based on the agent type:

### Reconnoiter Agent
- "Identifying target systems and services..."
- "Enumerating open ports and services..."
- "Mapping network topology and discovering hosts..."

### Vulnerability Scanner
- "Loading vulnerability signatures and patterns..."
- "Scanning for common vulnerabilities (CVEs)..."
- "Testing for authentication bypasses and injection flaws..."

### Exploiter
- "Analyzing identified vulnerabilities for exploitation..."
- "Selecting appropriate exploit techniques and payloads..."
- "Attempting controlled exploitation within RoE boundaries..."

### Reporter
- "Gathering findings from all assessment phases..."
- "Analyzing and correlating vulnerability data..."
- "Creating executive summary and technical details..."

## Architecture

### Key Components

1. **SubAgent Service** (`internal/subagent/activity.go`): Manages activity tracking and events
2. **Activity Page** (`internal/tui/page/activity.go`): TUI interface for monitoring
3. **Status Generator** (`internal/subagent/status_generator.go`): Generates intelligent status messages
4. **Agent Tool Integration** (`internal/agent/agenttool.go`): Hooks into subagent execution

### Event Flow

1. When orchestrator delegates a task, AgentTool starts tracking the activity
2. Activity Service publishes real-time events via pubsub system
3. Activity Page subscribes to events and updates the display
4. Users can abort tasks, which cancels the underlying agent execution

## Benefits

- **Transparency**: Users can see exactly what subagents are doing
- **Control**: Ability to abort tasks that seem problematic
- **Efficiency**: No need to wait for stuck or lengthy tasks
- **Trust**: Builds confidence through visibility into AI agent behavior
- **Debugging**: Helps identify issues with agent performance

## Technical Details

The system uses Go's context cancellation for safe task abortion and a pubsub event system for real-time updates. All activity data is stored in memory and automatically cleaned up after completion.