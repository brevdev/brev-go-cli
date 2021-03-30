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

	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
	"github.com/brevdev/brev-go-cli/internal/files"
)

func getSomeSetOfOptions(toComplete string) []string {
	return []string{"opt1", "opt2"}
}

func getEpNames() []string {
	var endpoints []brev_api.Endpoint
	files.ReadJSON(files.GetEndpointsPath(), &endpoints)

	var epNames []string
	for _, v := range endpoints {
		epNames = append(epNames, v.Name)
	}

	return epNames
}

func NewCmdEndpoint(context *cmdcontext.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "endpoint",
		Short: "Create, Run, or Remove Endpoints",
		Long: `Do any operation to your Brev endpoints. Ex:

		brev endpoint add NewEp
		brev endpoint run NewEp
		brev endpoint remove NewEp
	`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := cmdcontext.InvokeParentPersistentPreRun(cmd, args)
			if err != nil {
				return err
			}

			_, err = brev_api.CheckOutsideBrevErrorMessage(context)
			return err
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{}, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveDefault
		}}

	cmd.AddCommand(newCmdAdd(context))
	cmd.AddCommand(newCmdRemove(context))
	cmd.AddCommand(newCmdRun(context))
	// cmd.AddCommand(newCmdLog(context))
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return addEndpoint(name, context)
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeEndpoint(name, context)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "name of the endpoint")
	cmd.MarkFlagRequired("name")
	cmd.RegisterFlagCompletionFunc("name", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getEpNames(), cobra.ShellCompDirectiveNoSpace
	})
	return cmd
}

type Method int

const (
	GET Method = iota
	PUT
	POST
	DELETE
)

func newCmdRun(context *cmdcontext.Context) *cobra.Command {
	var name string
	var method string
	var arg []string
	var body string

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run your endpoint",
		Long: `Run your deployed endpoint in the console. Similar to cURL and Postman, etc
		ex:
			brev endpoint run MyEp
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEndpoint(name, method, arg, body, context)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "name of the endpoint")
	cmd.MarkFlagRequired("name")
	cmd.RegisterFlagCompletionFunc("name", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getEpNames(), cobra.ShellCompDirectiveNoSpace
	})
	cmd.Flags().StringVarP(&method, "method", "r", "GET", "http request method")
	cmd.MarkFlagRequired("method")
	cmd.RegisterFlagCompletionFunc("method", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"GET", "PUT", "POST", "DELETE"}, cobra.ShellCompDirectiveNoSpace
	})
	cmd.Flags().StringArrayVarP(&arg, "arg", "a", []string{}, "add query params")
	cmd.RegisterFlagCompletionFunc("arg", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoSpace
	})
	cmd.Flags().StringVarP(&body, "body", "b", "", "add json body")
	cmd.RegisterFlagCompletionFunc("body", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoSpace
	})

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
		RunE: func(cmd *cobra.Command, args []string) error {
			return logEndpoint(name, context)
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return listEndpoints(context)
		},
	}

	return cmd
}
