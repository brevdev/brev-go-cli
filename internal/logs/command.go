package logs

import (
	"github.com/spf13/cobra"

	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
	"github.com/brevdev/brev-go-cli/internal/terminal"
)

func NewCmdLogs(t *terminal.Terminal) *cobra.Command {
	var task string
	cmd := &cobra.Command{
		Use:     "logs",
		Short:   "Tail logs",
		Long:    `Tail logs of current project`,
		Example: `brev logs --task server`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := cmdcontext.InvokeParentPersistentPreRun(cmd, args)
			if err != nil {
				return err
			}

			_, err = brev_api.CheckOutsideBrevErrorMessage(t)
			return err
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{}, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveDefault
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return LogTask(task, t)

		},
	}

	cmd.Flags().StringVarP(&task, "task", "t", "", "The project task to tail (server, worker, packager, etc)")
	cmd.MarkFlagRequired("task")

	return cmd
}
