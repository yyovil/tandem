package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/yyovil/tandem/internal/models"
	"github.com/yyovil/tandem/internal/tools"
)

type Settings struct {
	AgentID      string           `json:"agentId,omitempty"`
	Description  string           `json:"description"`
	Goal         string           `json:"goal"`
	Instructions []string         `json:"instructions"`
	Model        models.ModelId   `json:"model,omitempty"`
	Name         string           `json:"name,omitempty"`
	Tools        []tools.ToolName `json:"tools"`
}

func NewSettings(path string) (settings Settings, err error) {

	// TODO: maybe we can look towards a specific path like .tandem/agent_name.json then.
	if path == "" {
		return settings, errors.New("empty string pass as path for settings")
	}

	file, err := os.Open(path)
	if err != nil {
		log.Println("Error opening settings file:", err)
		return settings, err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&settings)
	if err != nil {
		log.Println("Error decoding settings file:", err)
		return settings, err
	}

	return settings, nil
}

func (s Settings) GetSystemPrompt() string {
	return fmt.Sprintf(`
	You are an AI Agent that helps users with their tasks. You adhere strictly to the following guidelines:
	<description>%s</description>
	<your_goal>%s</your_goal>
	<instructions>%s</instructions>
	<additional_information>
	- use markdown for formatting your text response.
	- your name is %s. refer this name when users ask you about your identity. don't mention your underlying model.
	</additional_information>`,
		s.Description,
		s.Goal,
		s.Instructions,
		s.Name,
	)
}
