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
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/brevdev/brev-go-cli/internal/auth"
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
	"github.com/brevdev/brev-go-cli/internal/endpoint"
	"github.com/brevdev/brev-go-cli/internal/initialize"
	"github.com/brevdev/brev-go-cli/internal/package_project"
	"github.com/brevdev/brev-go-cli/internal/project"
	"github.com/brevdev/brev-go-cli/internal/version"
)

func main() {
	cmd := newCmdBrev()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func newCmdBrev() *cobra.Command {
	var verbose bool

	cmdContext := &cmdcontext.Context{}

	brevCommand := &cobra.Command{
		Use: "brev",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cmdContext.Init(verbose)

		},
	}

	brevCommand.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	createCmdTree(brevCommand, cmdContext)
	return brevCommand
}

func createCmdTree(brevCommand *cobra.Command, context *cmdcontext.Context) {
	brevCommand.AddCommand(endpoint.NewCmdEndpoint(context))
	brevCommand.AddCommand(auth.NewCmdLogin(context))
	brevCommand.AddCommand(package_project.NewCmdPackage(context))
	brevCommand.AddCommand(project.NewCmdProject(context))
	brevCommand.AddCommand(version.NewCmdVersion(context))
	brevCommand.AddCommand(initialize.NewCmdInit(context))
	brevCommand.AddCommand(&completionCmd)
}
