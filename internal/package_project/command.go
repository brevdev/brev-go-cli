package package_project

import (
	"github.com/spf13/cobra"

	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/brev_ctx"
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
	"github.com/brevdev/brev-go-cli/internal/terminal"
)

func getTopPyPiPackages() []string {
	return []string{"urllib3", "six", "boto3", "setuptools", "requests", "botocore", "idna", "certifi", "chardet", "pyyaml", "python-dateutil", "pip", "s3transfer", "wheel", "cffi", "rsa", "jmespath", "pyasn1", "numpy", "jinja"}
}

func GetPackages(t *terminal.Terminal) ([]brev_api.ProjectPackage, error) {
	brevCtx, err := brev_ctx.New()
	if err != nil {
		return nil, err
	}

	project, err := brevCtx.Local.GetProject()
	if err != nil {
		return nil, err
	}
	return brevCtx.Remote.GetPackages(*project, nil)
}

// This is just used for autocomplete, so failures can just return no autocompletions
func getCurrentPackages(t *terminal.Terminal) []string {
	packages, err := GetPackages(t)
	if err != nil {
		return []string{}
	}

	var packageNames []string
	for _, v := range packages {
		packageNames = append(packageNames, v.Name)
	}

	return packageNames
}

func NewCmdPackage(t *terminal.Terminal) *cobra.Command {
	cmd := &cobra.Command{
		Use:         "package",
		Annotations: map[string]string{"environment": ""},
		Short:       "Add or remove packages from your Brev project",
		Long:        "Add or remove python packages from your project (like pip)",
		Example: `  brev package add --name numpy
  brev package remove --name numpy`,
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
	cmd.AddCommand(newCmdList(t))

	return cmd
}

func newCmdAdd(t *terminal.Terminal) *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:     "add",
		Short:   "Add a python package to your project",
		Long:    "Installs a python package to your project (like pip).",
		Example: `  brev package add --name numpy`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return addPackage(name, t)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "name of the package")
	cmd.MarkFlagRequired("name")
	cmd.RegisterFlagCompletionFunc("name", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getTopPyPiPackages(), cobra.ShellCompDirectiveNoSpace
	})

	return cmd
}

func newCmdRemove(t *terminal.Terminal) *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:     "remove",
		Short:   "Remove a python package from your project",
		Long:    "Uninstalls a python package to your project (like pip).",
		Example: `  brev package remove --name numpy`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return removePackage(name, t)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "name of the package")
	cmd.MarkFlagRequired("name")
	cmd.RegisterFlagCompletionFunc("name", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getCurrentPackages(t), cobra.ShellCompDirectiveNoSpace
	})

	return cmd
}

func newCmdList(t *terminal.Terminal) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List installed packages",
		Long:    "List installed packages.",
		Example: `  brev package list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listPackages(t)
		},
	}

	return cmd
}
