package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/yaydraco/tandem/internal/agent"
	"github.com/yaydraco/tandem/internal/config"
	"github.com/yaydraco/tandem/internal/db"
	"github.com/yaydraco/tandem/internal/format"
	"github.com/yaydraco/tandem/internal/logging"
	"github.com/yaydraco/tandem/internal/message"
	"github.com/yaydraco/tandem/internal/session"
	"github.com/yaydraco/tandem/internal/tools"
)

type App struct {
	Sessions     session.Service
	Messages     message.Service
	Orchestrator agent.Service
	// ADHD: why we shouldn't initialise all the agents at once right in here? here's another thought. we don't want to have multiple agents of the same time, say couple of reconnoiters, doing some scanning because of the nature of the task in hand.
}

func New(ctx context.Context, conn *sql.DB) (*App, error) {
	q := db.New(conn)
	sessions := session.NewService(q)
	messages := message.NewService(q)

	app := &App{
		Sessions: sessions,
		Messages: messages,
	}

	var err error
	app.Orchestrator, err = agent.NewAgent(
		config.Orchestrator,
		app.Sessions,
		app.Messages,
		[]tools.BaseTool{agent.NewAgentTool(app.Sessions, app.Messages)},
	)

	if err != nil {
		logging.Error("Failed to create orchestrator agent", err)
		return nil, err
	}

	return app, nil
}

// RunNonInteractive handles the execution flow when a prompt is provided via CLI flag.
func (a *App) RunNonInteractive(ctx context.Context, prompt string, outputFormat string, quiet bool) error {
	logging.Info("Running in non-interactive mode")

	var spinner *format.Spinner
	if !quiet {
		spinner = format.NewSpinner("Thinking...")
		spinner.Start()
		defer spinner.Stop()
	}

	const maxPromptLengthForTitle = 100
	titlePrefix := "Non-interactive: "
	var titleSuffix string

	if len(prompt) > maxPromptLengthForTitle {
		titleSuffix = prompt[:maxPromptLengthForTitle] + "..."
	} else {
		titleSuffix = prompt
	}
	title := titlePrefix + titleSuffix

	sess, err := a.Sessions.Create(ctx, title)
	if err != nil {
		return fmt.Errorf("failed to create session for non-interactive mode: %w", err)
	}
	logging.Info("Created session for non-interactive run", "session_id", sess.ID)

	done, err := a.Orchestrator.Run(ctx, sess.ID, prompt)
	if err != nil {
		return fmt.Errorf("failed to start agent processing stream: %w", err)
	}

	result := <-done
	if result.Error != nil {
		if errors.Is(result.Error, context.Canceled) || errors.Is(result.Error, agent.ErrRequestCancelled) {
			logging.Info("Agent processing cancelled", "session_id", sess.ID)
			return nil
		}
		return fmt.Errorf("agent processing failed: %w", result.Error)
	}

	if !quiet && spinner != nil {
		spinner.Stop()
	}

	content := "No content available"
	if result.Message.Content().String() != "" {
		content = result.Message.Content().String()
	}

	fmt.Println(format.FormatOutput(content, outputFormat))

	logging.Info("Non-interactive run completed", "session_id", sess.ID)

	return nil
}