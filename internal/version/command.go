package version

import (
	"github.com/spf13/cobra"

	"github.com/brevdev/brev-go-cli/internal/terminal"
)

func NewCmdVersion(t *terminal.Terminal) *cobra.Command {
	cmd := &cobra.Command{
		Use: "version",
		RunE: func(cmd *cobra.Command, args []string) error {
			version, err := buildVersionString(t)
			if err != nil {
				t.Errprint(err, "Failed to determine version")
				return err
			}
			t.Vprint(version)
			return nil
		},
	}
	return cmd
}
