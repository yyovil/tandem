package login

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

// NewCommand returns `auth login`
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "log in to an auth provider",
		RunE: func(cmd *cobra.Command, args []string) error {
			// ensure config is loaded (best effort)
			wd, _ := os.Getwd()
			_, _ = config.Load(wd, false)

			provFlag, _ := cmd.Flags().GetString("provider")
			var selected authproviders.ProviderOption

			if provFlag == "" {
				// huh select for provider
				var selID string
				opts := []huh.Option[string]{}
				for _, p := range authproviders.Providers {
					opts = append(opts, huh.NewOption(p.Label, p.ID))
				}
				form := huh.NewForm(
					huh.NewGroup(
						huh.NewSelect[string]().
							Title("select a provider to log in").
							Options(opts...).
							Value(&selID),
					),
				)
				if err := form.Run(); err != nil {
					return err
				}
				for _, p := range authproviders.Providers {
					if p.ID == selID {
						selected = p
						break
					}
				}
				if selected.ID == "" {
					return fmt.Errorf("no provider selected")
				}
			} else {
				for _, p := range authproviders.Providers {
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
				// try auto-discover token
				if token, err := config.LoadGitHubToken(); err == nil && strings.TrimSpace(token) != "" {
					cfg := config.Get()
					if cfg == nil {
						return fmt.Errorf("config not loaded")
					}
					if cfg.Providers == nil {
						cfg.Providers = make(map[models.ModelProvider]config.Provider)
					}
					p := cfg.Providers[models.ProviderCopilot]
					p.APIKey = token
					p.Disabled = false
					cfg.Providers[models.ProviderCopilot] = p
					fmt.Fprintln(cmd.OutOrStdout(), "Logged in to GitHub Copilot.")
					return nil
				}

				// fallback: prompt for token
				var token string
				form := huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("enter your GitHub Copilot token").
							Placeholder("paste token").
							Password(true).
							Value(&token),
					),
				)
				if err := form.Run(); err != nil {
					return err
				}
				token = strings.TrimSpace(token)
				if token == "" {
					return fmt.Errorf("no token provided")
				}
				cfg := config.Get()
				if cfg == nil {
					return fmt.Errorf("config not loaded")
				}
				if cfg.Providers == nil {
					cfg.Providers = make(map[models.ModelProvider]config.Provider)
				}
				pp := cfg.Providers[models.ProviderCopilot]
				pp.APIKey = token
				pp.Disabled = false
				cfg.Providers[models.ProviderCopilot] = pp
				fmt.Fprintln(cmd.OutOrStdout(), "Logged in to GitHub Copilot.")
				return nil
			default:
				return fmt.Errorf("provider not yet supported: %s", selected.Label)
			}
		},
	}

	cmd.Flags().StringP("provider", "p", "", "provider id or name (e.g., github-copilot)")
	return cmd
}
