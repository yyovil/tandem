package auth

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/yyovil/tandem/internal/cmd/auth/providers"
)

// NewListCommand returns the `auth list` subcommand
func NewListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list available auth providers",
		RunE: func(cmd *cobra.Command, args []string) error {
			// show providers using huh select (display-only)
			var sel string
			options := []huh.Option[string]{}
			for _, p := range providers.Providers {
				options = append(options, huh.NewOption(fmt.Sprintf("%s (%s)", p.Label, p.ID), p.ID))
			}
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title("oauth providers available").
						Options(options...).
						Value(&sel),
				),
			)
			return form.Run()
		},
	}
}
