package auth

import (
	"github.com/spf13/cobra"

	"github.com/brevdev/brev-go-cli/internal/terminal"
)

func NewCmdLogin(t *terminal.Terminal) *cobra.Command {
	cmd := &cobra.Command{
		Use: "login",
		RunE: func(cmd *cobra.Command, args []string) error {
			return login(t)
		},
	}
	return cmd
}
