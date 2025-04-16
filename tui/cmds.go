package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
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
}

// find a better way to provide the available models option.
const (
	ModelGeminiFlashLite string = "gemini-2.0-flash-lite"
)

// this cmd makes the HTTP POST request to an endpoint to receive the response stream.
func (m *Model) GetCompletionStreamCmd() tea.Cmd {
	return func() tea.Msg {

		reqBody := RunRequest{
			SessionId: m.sessionId,
			Model:     m.model,
			UserId:    m.userId,
			Message:   m.message,
		}

		serialisedReqBody, err := json.Marshal(reqBody)
		if err != nil {
			// also notify why serialization failed in debug mode.
			return m
		}

		sse, err := http.Post(os.Getenv("ENDPOINT"), "application/json", bytes.NewBuffer(serialisedReqBody))

		if err != nil {
			// fmt.Println("POST req failed:", err)
			return m
		}

		defer sse.Body.Close()

		if sse.StatusCode != http.StatusOK {
			// fmt.Println("Stream request failed:", sse.Status)
			return m
		}

		streamReader := bufio.NewReader(sse.Body)

		for {
			line, err := streamReader.ReadString('\n')
			if err != nil {
				// exhausted.
				if err == io.EOF {
					// we have to return the complete streamed response.
					break
				}
				// fmt.Println("Error reading stream:", err)
				return m
			}
			m.streamingResponse += line
		}

		return m
	}
}
