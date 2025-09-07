package auth

import (
	"github.com/spf13/cobra"
	"github.com/yyovil/tandem/internal/cmd/auth/login"
	"github.com/yyovil/tandem/internal/cmd/auth/logout"
)

// NewCommand creates the root `auth` command and wires subcommands
func NewCommand() *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "authorize tandem to use 3rd-party providers",
	}

	// attach subcommands
	authCmd.AddCommand(login.NewCommand())
	authCmd.AddCommand(logout.NewCommand())
	authCmd.AddCommand(NewListCommand())

	return authCmd
}
