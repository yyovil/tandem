package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// Request model for running an agent
type RunRequest struct {
	SessionId string `json:"session_id"`
	Model     string `json:"model"`
	UserId    string `json:"user_id"`
	Message   string `json:"message"`
	Stream    bool   `json:"stream"`
}

// find a better way to provide the available models option.
const (
	ModelGemini20FlashLite string = "gemini-2.0-flash-lite"
	ModelGemini25ProPreview0325 string = "gemini-2.5-pro-preview-03-25"
	ModelGemini25FlashPreview0417 string = "gemini-2.5-flash-preview-04-17"
)

// this cmd makes the HTTP POST request to an endpoint to receive the response stream.
func (m *Model) GetCompletionStreamCmd() tea.Cmd {
	return func() tea.Msg {
		reqBody := RunRequest{
			SessionId: m.sessionId,
			Model:     ModelGemini25FlashPreview0417,
			UserId:    m.userId,
			Message:   m.message,
			Stream:    true,
		}

		serialisedReqBody, err := json.Marshal(reqBody)
		if err != nil {
			// also notify why serialization failed in debug mode.
			return m
		}

		endpoint := os.Getenv("ENDPOINT")
		if endpoint == "" {
			log.Println("ENDPOINT not set")
			return m
		}

		sse, err := http.Post(os.Getenv("ENDPOINT"), "application/json", bytes.NewBuffer(serialisedReqBody))

		if err != nil {
			log.Println("POST req failed:", sse.Status)
			return m
		}

		defer sse.Body.Close()
		if sse.StatusCode != http.StatusOK {
			log.Println("Stream request failed:", sse.Status)
			return m
		}

		streamReader := bufio.NewReader(sse.Body)

		line, err := streamReader.ReadString('\n')
		if err == io.EOF {
			return streamDoneMsg{}
		}
		return streamChunkMsg{Text: line}
	}
}
