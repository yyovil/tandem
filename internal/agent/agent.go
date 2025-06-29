package agent

import (
	"context"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yyovil/tandem/internal/settings"
	"github.com/yyovil/tandem/internal/tools"
)

type Provider interface {
	// returns a stream of messages based on the provided context, chat history, and agent settings.
	GetStream(ctx context.Context, messages []Message, settings settings.Settings) <-chan Message

	// returns messages in provider specific api from tandem specified api format.
	FromMessages(messages []Message) any

	// returns message from provider specific api to tandem specified api format.
	ToMessage(message any) Message

	// returns tools in provider specific api format.
	GetToolsForProvider(tools []tools.ToolName) any

	// returns parameters in provider specific api from tandem specified format.
	FromParameters(params tools.ToolParameters) any

	// returns schema for the given tool parameter in provider specific api format.
	ToSchema(param tools.Param) any
}

// NOTE: Use this method to get a new provider for the agent. use the user provided settings and fallback to some defaults.
// func NewProvider() (Provider, error) {}

type Agent struct {
	// provider behind the agent powering its intelligence.
	Provider Provider
	// agent settings that you could provide using a json file. here's the schema for it: https://github.com/yyovil/tandem/blob/fuck_agno/agents/agent.schema.json
	Settings *settings.Settings
}

// Executes the agent's logic, generating content based on the chat history and user settings
func (a Agent) Run(ctx context.Context, history []Message) tea.Cmd {
	/*
		!TODO:
		1. check if there's a chat session available already. (waiting for session persistence)
	*/
	return func() tea.Msg {
		return StreamCreated{
			Stream: a.Provider.GetStream(ctx, history, *a.Settings),
		}
	}
}

func NewAgent(settings settings.Settings) (Agent, error) {
	// provider, err := NewProvider()
	// if err != nil {
	// 	return Agent{}, err
	// }

	return Agent{
		Settings: &settings,
		// Provider: provider,
	}, nil
}
