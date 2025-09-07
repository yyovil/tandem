package logout

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	authproviders "github.com/yyovil/tandem/internal/cmd/auth/providers"
	"github.com/yyovil/tandem/internal/config"
	"github.com/yyovil/tandem/internal/models"
)

// NewCommand returns `auth logout`
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "log out from an auth provider",
		RunE: func(cmd *cobra.Command, args []string) error {
			// ensure config is loaded
			wd, _ := os.Getwd()
			_, _ = config.Load(wd, false)

			provFlag, _ := cmd.Flags().GetString("provider")

			// build list of logged-in providers (for now only copilot)
			loggedIn := []authproviders.ProviderOption{}
			if cfg := config.Get(); cfg != nil {
				if p, ok := cfg.Providers[models.ProviderCopilot]; ok && strings.TrimSpace(p.APIKey) != "" && !p.Disabled {
					loggedIn = append(loggedIn, authproviders.ProviderOption{ID: "github-copilot", Label: "GitHub Copilot", Prov: models.ProviderCopilot})
				}
			}
			if len(loggedIn) == 0 {
				// allow explicit selection via flag anyway
				loggedIn = append(loggedIn, authproviders.ProviderOption{ID: "github-copilot", Label: "GitHub Copilot", Prov: models.ProviderCopilot})
			}

			var selected authproviders.ProviderOption
			if provFlag == "" {
				var selID string
				opts := []huh.Option[string]{}
				for _, p := range loggedIn {
					opts = append(opts, huh.NewOption(p.Label, p.ID))
				}
				form := huh.NewForm(
					huh.NewGroup(
						huh.NewSelect[string]().
							Title("select a provider to log out").
							Options(opts...).
							Value(&selID),
					),
				)
				if err := form.Run(); err != nil {
					return err
				}
				for _, p := range loggedIn {
					if p.ID == selID {
						selected = p
						break
					}
				}
				if selected.ID == "" {
					return fmt.Errorf("no provider selected")
				}
			} else {
				for _, p := range loggedIn {
					if p.ID == provFlag || strings.EqualFold(p.Label, provFlag) {
						selected = p
						break
					}
				}
				if selected.ID == "" {
					return fmt.Errorf("unknown provider: %s", provFlag)
				}
			}

			switch selected.Prov {
			case models.ProviderCopilot:
				cfg := config.Get()
				if cfg == nil {
					return fmt.Errorf("config not loaded")
				}
				if cfg.Providers == nil {
					cfg.Providers = make(map[models.ModelProvider]config.Provider)
				}
				p := cfg.Providers[models.ProviderCopilot]
				p.APIKey = ""
				p.Disabled = true
				cfg.Providers[models.ProviderCopilot] = p
				fmt.Fprintln(cmd.OutOrStdout(), "Logged out from GitHub Copilot.")
				return nil
			default:
				return fmt.Errorf("provider not yet supported: %s", selected.Label)
			}
		},
	}

	cmd.Flags().StringP("provider", "p", "", "provider id or name (e.g., github-copilot)")
	return cmd
}
