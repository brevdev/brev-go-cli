package sync

import (
	"github.com/spf13/cobra"

	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
	"github.com/brevdev/brev-go-cli/internal/terminal"
)

func NewCmdPush(t *terminal.Terminal) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push",
		Short: "Push your local changes to remote",
		Long: `To push your local changes:

			brev push
		`, PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := cmdcontext.InvokeParentPersistentPreRun(cmd, args)
			if err != nil {
				return err
			}

			_, err = brev_api.CheckOutsideBrevErrorMessage(t)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return push(t)
		},
	}

	return cmd
}

func NewCmdPull(t *terminal.Terminal) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull latest changes from your server",
		Long: `To pull latest changes:

			brev pull
		`, PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := cmdcontext.InvokeParentPersistentPreRun(cmd, args)
			if err != nil {
				return err
			}

			_, err = brev_api.CheckOutsideBrevErrorMessage(t)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return pull(t)
		},
	}

	return cmd
}
