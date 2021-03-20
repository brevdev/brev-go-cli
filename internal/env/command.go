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
package env

import (
	"fmt"

	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
	"github.com/spf13/cobra"
)


func NewCmdEnv(context *cmdcontext.Context) *cobra.Command {
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

		_, err = brev_api.CheckOutsideBrevErrorMessage(context)
		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("env called")
	},
}
	
	cmd.RemoveCommand(newCmdRemove(context))
	cmd.AddCommand(newCmdAdd(context))
	cmd.AddCommand(newCmdRemove(context))


	return cmd
}

func newCmdAdd(context *cmdcontext.Context) *cobra.Command {
	var name string
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add an encrypted environment variable",
		Long: `To add an environment variable:

			brev env add --name XYZ

		You will then be prompted for the value.
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return addVariable(name, context)
		},
	}
	cmd.Flags().StringVarP(&name, "name", "n", "", "variable name")
	cmd.MarkFlagRequired("name")

	return cmd
}

func newCmdRemove(context *cmdcontext.Context) *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove an environment variable",
		Long: `To remove an environment variable:

			brev env remove --name XYZ
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeVariable(name, context)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "variable name")
	cmd.MarkFlagRequired("name")

	return cmd
}
