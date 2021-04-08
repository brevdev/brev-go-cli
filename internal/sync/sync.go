package sync

import (
	"github.com/spf13/cobra"

	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
	"github.com/brevdev/brev-go-cli/internal/terminal"
)

func NewCmdPush(t *terminal.Terminal) *cobra.Command {
	cmd := &cobra.Command{
		Use:         "push",
		Annotations: map[string]string{"code": ""},
		Short:       "Push your local changes to remote",
		Example:     `  brev push`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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
		Use:         "pull",
		Annotations: map[string]string{"code": ""},
		Short:       "Pull latest changes from your server",
		Example:     `  brev pull`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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

func NewCmdDiff(t *terminal.Terminal) *cobra.Command {

	cmd := &cobra.Command{
		Use:         "diff",
		Annotations: map[string]string{"code": ""},
		Short:       "See a diff of your local changes compared to what's deployed in the console",
		Long: `To see a diff of your local changes compared to what's deployed in the console,
			from an active brev project directory, run:

			brev diff
		`, PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := cmdcontext.InvokeParentPersistentPreRun(cmd, args)
			if err != nil {
				return err
			}

			_, err = brev_api.CheckOutsideBrevErrorMessage(t)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return diffCmd(t)
		},
	}

	return cmd
}
