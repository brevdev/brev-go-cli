/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package status

import (
	"fmt"

	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/brev_ctx"
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
	"github.com/spf13/cobra"
)

func NewCmdStatus(context *cmdcontext.Context) *cobra.Command {

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

			_, err = brev_api.CheckOutsideBrevErrorMessage(context)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			status(context)
			return nil
		},
	}

	return cmd
}

func status(context *cmdcontext.Context) error {

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

	fmt.Fprintf(context.VerboseOut, "\nProject %s", project.Name)

	// Print package info
	if len(packages) == 0 {
		fmt.Fprintln(context.VerboseOut, "\n\tNo packages installed.")
	} else {
		fmt.Fprintln(context.VerboseOut, "\n\tPackages:")
	}

	for _, v := range packages {
		fmt.Fprintf(context.VerboseOut, "\t\t %s==%s %s\n", v.Name, v.Version, v.Status)
	}

	// Print Endpoint info
	if len(endpoints) == 0 {
		fmt.Fprintln(context.VerboseOut, "\nYour project doesn't have any endpoints. Try running \n \t\t brev endpoint add --name newEP")
	} else {
		fmt.Fprintln(context.VerboseOut, "\n\tEndpoints:")

		for _, v := range endpoints {
			fmt.Fprintf(context.VerboseOut, "\n\t\t%s\n\t\t\t%s", v.Name, project.Domain+v.Uri)
		}
	}

	return nil
}
