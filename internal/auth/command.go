package auth

import (
	"github.com/spf13/cobra"

	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
)

func NewCmdLogin(context *cmdcontext.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "login",
		RunE: func(cmd *cobra.Command, args []string) error {
			return login(context)
		},
	}
	return cmd
}
