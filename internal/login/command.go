package login

import (
	"github.com/spf13/cobra"

	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
)

func NewCmdLogin(context *cmdcontext.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "login",
		Run: func(cmd *cobra.Command, args []string) {
			authenticateWithCotter()
		},
	}
	return cmd
}
