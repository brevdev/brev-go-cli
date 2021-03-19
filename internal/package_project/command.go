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
package package_project

import (
	"github.com/spf13/cobra"

	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
)

func getTopPyPiPackages() []string {
	return []string{"urllib3","six","boto3","setuptools","requests","botocore","idna","certifi","chardet","pyyaml","python-dateutil","pip","s3transfer","wheel","cffi","rsa","jmespath","pyasn1","numpy","jinja",}
}

func NewCmdPackage(context *cmdcontext.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "package",
		Short: "Add or remove packages from your Brev project",
		Long: `Add or remove python packages from your project (like pip):
		ex:
			brev package add --name numpy
			brev package remove --name numpy
	`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := cmdcontext.InvokeParentPersistentPreRun(cmd, args)
			if err != nil {
				return err
			}

			_, err = brev_api.CheckOutsideBrevErrorMessage(context)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return logic(context)
		},
	}

	cmd.AddCommand(newCmdAdd(context))

	return cmd
}

func newCmdAdd(context *cmdcontext.Context) *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a python package to your project",
		Long: `Installs a python package to your project (like pip)
			ex: 
				brev package add --name numpy
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return addPackage(name, context)
		},

	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "name of the package")
	cmd.MarkFlagRequired("name")
	cmd.RegisterFlagCompletionFunc("name", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getTopPyPiPackages(), cobra.ShellCompDirectiveNoSpace
	})
	

	return cmd
}
