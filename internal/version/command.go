package version

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/brevdev/brev-go-cli/internal/auth"
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
)

func NewCmdVersion(context *cmdcontext.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "version",
		RunE: func(cmd *cobra.Command, args []string) error {
			version, err := buildVersionString(context)
			if err != nil {
				return err
			}
			fmt.Fprintln(context.VerboseOut, version)

			token, _ := auth.GetToken()
			fmt.Println(token)
			return nil
		},
	}
	return cmd
}
