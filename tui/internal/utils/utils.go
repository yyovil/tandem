package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type Model string

const (
	GEMINI_2_5_FLASH_PREVIEW_04_17    Model = "gemini-2.5-flash-preview-04-17"
	GEMINI_2_5_PRO_EXPERIMENTAL_03_25 Model = "gemini-2.5-pro-exp-03-25"
	GEMINI_2_0_FLASH                  Model = "gemini-2.0-flash"
	GEMINI_2_0_FLASH_LITE             Model = "gemini-2.0-flash-lite"
)

type RunRequest struct {
	Prompt    string `json:"message"`
	Stream    string `json:"stream"`
	Model     Model  `json:"model"`
	UserId    string `json:"user_id"`
	SessionId string `json:"session_id"`
	// Attachments [][]byte `json:"attachments"`
	// TODO: support attachments later
}

func Post(Prompt string) (*http.Request, error) {
	runRequestBody := RunRequest{
		Prompt:    Prompt,
		Stream:    "true",
		Model:     GEMINI_2_5_FLASH_PREVIEW_04_17,
		UserId:    "tanishq",
		SessionId: "dummy-pentest-session",
	}

	jsonData, err := json.Marshal(runRequestBody)
	if err != nil {
		log.Println("error marshaling JSON:", err.Error())
		return nil, err
	}

	agentId := "Mr. Burnham"

	runRequest, err := http.NewRequest(http.MethodPost, "http://localhost:8000/v1/agents/"+agentId+"/runs", bytes.NewBuffer(jsonData))

	if err != nil {
		log.Println("error while creating a new request.", err.Error())
		// ADHD: think what happens when we can't create a request?
		return nil, err
	}

	return runRequest, nil
}

/*
TODO: we need a util that calculates the appropriate width of a cmp that when set| upon window resize, doesn't overflows the content inside of it.
*/
