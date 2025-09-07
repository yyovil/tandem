package providers

import "github.com/yyovil/tandem/internal/models"

// provideroption represents an auth provider
type ProviderOption struct {
    ID    string
    Label string
    Prov  models.ModelProvider
}

// Providers currently available for oauth
var Providers = []ProviderOption{
    {ID: "github-copilot", Label: "GitHub Copilot", Prov: models.ProviderCopilot},
}
