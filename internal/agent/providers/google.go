package providers

import (
	"context"
	"log"
	"os"

	"google.golang.org/genai"
)

type GoogleProvider struct {
	Client *genai.Client
}

func NewGoogleProvider() (GoogleProvider, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})

	if err != nil {
		log.Println("Error creating Gemini client:", err.Error())
		return GoogleProvider{}, err
	}

	return GoogleProvider{
		Client: client,
	}, nil
}
