package agent

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/yyovil/tandem/internal/agent/providers"
	"github.com/yyovil/tandem/internal/agent/tools"
)

type Settings struct {
	AgentID      string             `json:"agentId,omitempty"`
	Description  string             `json:"description"`
	Goal         string             `json:"goal"`
	Instructions []string           `json:"instructions"`
	Model        *providers.ModelId `json:"model,omitempty"`
	Name         string             `json:"name,omitempty"`
	Tools        []tools.ToolName   `json:"tools"`
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