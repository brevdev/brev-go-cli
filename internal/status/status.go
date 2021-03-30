
package status

import (
	"github.com/spf13/cobra"

	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/brev_ctx"
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
	"github.com/brevdev/brev-go-cli/internal/terminal"
)

func NewCmdStatus(t *terminal.Terminal) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Get the latest project metadata",
		Long: `See high level on your project. Ex:

			brev status

		`, PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := cmdcontext.InvokeParentPersistentPreRun(cmd, args)
			if err != nil {
				return err
			}

			_, err = brev_api.CheckOutsideBrevErrorMessage(t)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			status(t)
			return nil
		},
	}

	return cmd
}

func status(t *terminal.Terminal) error {

	brevCtx, err := brev_ctx.New()
	if err != nil {
		return err
	}

	project, err := brevCtx.Local.GetProject()
	if err != nil {
		return err
	}

	endpoints, err := brevCtx.Local.GetEndpoints(&brev_ctx.GetEndpointsOptions{
		ProjectID: project.Id,
	})
	if err != nil {
		return err
	}

	packages, err := brevCtx.Remote.GetPackages(*project, nil)
	if err != nil {
		return err
	}

	t.Vprint(t.Yellow("\nProject %s", project.Name))

	// Print package info
	if len(packages) == 0 {
		t.Vprint("\n\tNo packages installed.")
	} else {
		t.Vprint(t.Yellow("\n\tPackages:"))
	}

	for _, v := range packages {
		if v.Status == "pending" {
			t.Vprintf("\t\t%s==%s %s\n", v.Name, v.Version, t.Yellow(v.Status))
		} else if v.Status == "installed" {
			t.Vprintf("\t\t%s==%s %s\n", v.Name, v.Version, t.Green(v.Status))
		} else {
			t.Vprintf("\t\t%s==%s %s\n", v.Name, v.Version, t.Red(v.Status))
		}
	}

	// Print Endpoint info
	if len(endpoints) == 0 {
		t.Vprint("\nYour project doesn't have any endpoints. Try running \n \t\t brev endpoint add --name newEP")
	} else {
		t.Vprint(t.Yellow("\n\tEndpoints:"))

		for _, v := range endpoints {
			t.Vprint("\n\t\t")
			t.Vprint(t.Yellow("%s", v.Name))
			t.Vprint("\n\t\t\t")
			t.Vprintf("%s", project.Domain+v.Uri)
		}
	}

	return nil
}
