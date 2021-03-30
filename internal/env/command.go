
package env

import (
	"github.com/spf13/cobra"

	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/brev_ctx"
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
	"github.com/brevdev/brev-go-cli/internal/terminal"
)

func NewCmdEnv(t *terminal.Terminal) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env",
		Short: "Add or remove environment variables",
		Long: `Use the Brev secrets manager for encrypted variables that get used at runtime.
		
		ex: 
			brev env add XYZ

		code usage:
			import variables
			...
			print(variables.XYZ)
	`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := cmdcontext.InvokeParentPersistentPreRun(cmd, args)
			if err != nil {
				return err
			}

			_, err = brev_api.CheckOutsideBrevErrorMessage(t)
			return err
		},
	}

	cmd.AddCommand(newCmdAdd(t))
	cmd.AddCommand(newCmdRemove(t))

	return cmd
}

func newCmdAdd(t *terminal.Terminal) *cobra.Command {
	var name string
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add an encrypted environment variable",
		Long: `To add an environment variable:

			brev env add --name XYZ

		You will then be prompted for the value.
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return addVariable(name, t)
		},
	}
	cmd.Flags().StringVarP(&name, "name", "n", "", "variable name")
	cmd.MarkFlagRequired("name")

	return cmd
}

func newCmdRemove(t *terminal.Terminal) *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove an environment variable",
		Long: `To remove an environment variable:

			brev env remove --name XYZ
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeVariable(name, t)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "variable name")
	cmd.MarkFlagRequired("name")
	cmd.RegisterFlagCompletionFunc("name", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVariables(), cobra.ShellCompDirectiveNoSpace
	})

	return cmd
}

// For shell completions, let the command raise an error
// if something fails here, just return nil
// i.e. don't provide completion but let user continue
func getVariables() []string {
	brevCtx, err := brev_ctx.New()
	if err != nil {
		return nil
	}

	project, err := brevCtx.Local.GetProject()
	if err != nil {
		return nil
	}

	vars, err := brevCtx.Remote.GetVariables(*project, nil)
	if err != nil {
		return nil
	}

	var varNames []string
	for _, v := range vars {
		varNames = append(varNames, v.Name)
	}
	return varNames

}
