package providers

import (
	"context"

	"github.com/yyovil/tandem/internal/agent/tools"
	"github.com/yyovil/tandem/internal/chat"
)

type Provider interface {
	// GetStream returns a stream of messages based on the provided context, chat history, and tools.
	GetStream(ctx context.Context, history chat.History, tools []tools.Tool) Stream
	ConvertMessages()
}

// NOTE: Use this method to get a new provider for the agent.
func NewProvider() (Provider, error) {

}
