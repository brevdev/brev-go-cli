package version

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
)

const command = "version"

func NewCmdVersion(context *cmdcontext.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: command,
		RunE: func(cmd *cobra.Command, args []string) error {
			version, err := buildVersionString(context)
			if err != nil {
				return err
			}
			fmt.Fprintln(context.VerboseOut, version)
			return nil
		},
	}
	return cmd
}
