package providers

import (
	"context"
	"log"
	"os"

	"google.golang.org/genai"
)

type GeminiProvider struct {
	Client *genai.Client
}

func NewGeminiProvider() (GeminiProvider, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})

	if err != nil {
		log.Println("Error creating Gemini client:", err.Error())
		return GeminiProvider{}, err
	}

	return GeminiProvider{
		Client: client,
	}, nil
}

