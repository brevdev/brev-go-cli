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
package endpoint

import (
	"github.com/spf13/cobra"

	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
)

func NewCmdEndpoint(context *cmdcontext.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "endpoint",
		Short: "Create, Run, or Remove Endpoints",
		Long: `Do any operation to your Brev endpoints. Ex:

		brev endpoint add NewEp
		brev endpoint run NewEp
		brev endpoint remove NewEp
	`,
	}

	cmd.AddCommand(newCmdAdd(context))
	cmd.AddCommand(newCmdRemove(context))
	cmd.AddCommand(newCmdRun(context))
	cmd.AddCommand(newCmdLog(context))
	cmd.AddCommand(newCmdList(context))

	return cmd
}

func newCmdAdd(context *cmdcontext.Context) *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add an endpoint to your project.",
		Long: `Add an endpoint to your project. This will also create the file in your directory.
		ex:
			brev endpoint add NewEp
		`,
		Run: func(cmd *cobra.Command, args []string) {
			add_endpoint(name)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "name of the endpoint")

	return cmd
}

func newCmdRemove(context *cmdcontext.Context) *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove an endpoint from your project.",
		Long: `Remove an endpoint from your project. This will also remove the file from your directory.
		ex:
			brev endpoint remove NewEp
		`,
		Run: func(cmd *cobra.Command, args []string) {
			remove_endpoint(name)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "name of the endpoint")

	return cmd
}

func newCmdRun(context *cmdcontext.Context) *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run your endpoint",
		Long: `Run your endpoint  on the remote server. Similar to cURL and Postman, etc
		ex:
			brev endpoint run MyEp
		`,
		Run: func(cmd *cobra.Command, args []string) {
			run_endpoint(name)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "name of the endpoint")

	return cmd
}

func newCmdLog(context *cmdcontext.Context) *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "log",
		Short: "Log an endpoint.",
		Long: `Get logs for any endpoint.
		ex:
			brev endpoint log NewEp
		`,
		Run: func(cmd *cobra.Command, args []string) {
			log_endpoint(name)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "name of the endpoint")

	return cmd
}

func newCmdList(context *cmdcontext.Context) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List endpoints in your project.",
		Long: `List endpoints in your project. This will print your URLs.
		ex:
			brev endpoint list
		`,
		Run: func(cmd *cobra.Command, args []string) {
			list_endpoints()
		},
	}

	return cmd
}
