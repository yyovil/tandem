package agent

import (
	"context"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yyovil/tandem/internal/agent/providers"
	"github.com/yyovil/tandem/internal/chat"
	"google.golang.org/genai"
)

type Service interface {
	Run()
}

type Agent struct {
	Id       string
	Provider providers.Provider
	Chat     chat.Chat //I'm not sure what's the use of this.
	Settings Settings
}

// Run executes the agent's logic, generating content based on the chat history and settings.
func (a Agent) Run(ctx context.Context, content chan<- genai.Content) tea.Cmd {
	/*
		TODO:
		1. check if there's a chat session available already.
	*/

	return func() tea.Msg {
		stream := a.Provider.GetStream(ctx, a.Chat.History, a.Settings.Tools)

		return nil
	}
}

func NewAgent(settings Settings) (Agent, error) {
	provider, err := providers.NewProvider()
	if err != nil {
		return Agent{}, err
	}

	file, err := os.Open("path/to/file.json")
	if err != nil {
		// handle error
		log.Println("Error opening file:", err.Error())
	}
	defer file.Close()

	return Agent{
		Settings: Settings{},
		Provider: provider,
	}, nil
}

/*
TODO:
1. Take the config from the internal/agents and build the agents from the definition.
*/
