package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/charmbracelet/x/ansi"
)

type Model string

type Attachment struct {
	Url      string `json:"url"`
	Filepath string `json:"filepath"`
	Content  string `json:"content"`
	MimeType string `json:"mime_type"`
}

const (
	GEMINI_2_5_FLASH_PREVIEW_04_17    Model = "gemini-2.5-flash-preview-04-17"
	GEMINI_2_5_PRO_EXPERIMENTAL_03_25 Model = "gemini-2.5-pro-exp-03-25"
	GEMINI_2_0_FLASH                  Model = "gemini-2.0-flash"
	GEMINI_2_0_FLASH_LITE             Model = "gemini-2.0-flash-lite"
)

type RunRequest struct {
	Prompt      string       `json:"message"`
	Stream      string       `json:"stream"`
	Model       Model        `json:"model"`
	UserId      string       `json:"user_id"`
	SessionId   string       `json:"session_id"`
	Attachments []Attachment `json:"attachments"`
}

func GetPostRequest(prompt string, attachments []Attachment) (*http.Request, error) {

	runRequestBody := RunRequest{
		Prompt:      prompt,
		Stream:      "true",
		Model:       GEMINI_2_5_FLASH_PREVIEW_04_17,
		UserId:      "tanishq",
		SessionId:   "dummy-pentest-session-2",
		Attachments: nil,
	}

	if len(attachments) > 0 {
		for _, attachment := range attachments {
			runRequestBody.Attachments = append(runRequestBody.Attachments, Attachment{
				Filepath: attachment.Filepath,
				MimeType: attachment.MimeType,
				// Content:  attachment.Content,
			})
		}
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

func Wordwrap(content string, width int) string {
	// ADHD: these breakpoints are silly.
	var breakpoints string = " ,-"
	return ansi.Wordwrap(content, width, breakpoints)
}

// Clamp ensures that the value is within the specified min and max range.
func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
