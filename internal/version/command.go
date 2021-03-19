package version

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
)

func NewCmdVersion(context *cmdcontext.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "version",
		RunE: func(cmd *cobra.Command, args []string) error {
			version, err := buildVersionString(context)
			if err != nil {
				context.PrintErr("Failed to determine version", err)
				return err
			}
			fmt.Fprintln(context.VerboseOut, version)
			return nil
		},
	}
	return cmd
}
