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

func getSomeSetOfOptions(toComplete string) []string {
	return []string{"opt1", "opt2"}
}

func getEpNames() []string {
	print("hhhhhhhh")
	return []string{"ep1", "ep2", "ep3"}
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
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			
			if (args[len(args)-1] == "run" || args[len(args)-1] == "remove") {
				return getEpNames(), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveDefault
			}

			return getSomeSetOfOptions(toComplete), cobra.ShellCompDirectiveNoFileComp
		}}

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
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			// fmt.Println("heeyy from this function")
			return getEpNames(), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveDefault
			// if len(args) != 0 {
			// 	return nil, cobra.ShellCompDirectiveNoFileComp
			// }
			
			// if (args[len(args)-1] == "run" || args[len(args)-1] == "remove") {
			// 	return getEpNames(), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveDefault
			// }

			// return getSomeSetOfOptions(toComplete), cobra.ShellCompDirectiveNoFileComp
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "name of the endpoint")

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
		Long: `Run your endpoint  on the remote server. Similar to cURL and Postman, etc
		ex:
			brev endpoint run MyEp
		`,
		Run: func(cmd *cobra.Command, args []string) {
			run_endpoint(name, method, arg, body)
			// for _, v := range arg {
			// 	fmt.Println(v)
			// }
		},

	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "name of the endpoint")
	cmd.MarkFlagRequired("name")
	cmd.RegisterFlagCompletionFunc("name", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getEpNames(), cobra.ShellCompDirectiveNoSpace})
	cmd.Flags().StringVarP(&method, "method", "r", "GET", "http request method")
	cmd.MarkFlagRequired("method")
	cmd.RegisterFlagCompletionFunc("method", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"GET", "PUT", "POST", "DELETE"}, cobra.ShellCompDirectiveNoSpace})
	cmd.Flags().StringArrayVarP(&arg, "arg", "a", []string{}, "add query params")
	cmd.Flags().StringVarP(&body, "body", "b", "", "add json body")

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
