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
	"github.com/brevdev/brev-go-cli/internal/brev_errors"
	"os"

	"github.com/spf13/cobra"

	"github.com/brevdev/brev-go-cli/internal/auth"
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
	"github.com/brevdev/brev-go-cli/internal/endpoint"
	"github.com/brevdev/brev-go-cli/internal/env"
	"github.com/brevdev/brev-go-cli/internal/initialize"
	"github.com/brevdev/brev-go-cli/internal/package_project"
	"github.com/brevdev/brev-go-cli/internal/status"
	"github.com/brevdev/brev-go-cli/internal/sync"
	"github.com/brevdev/brev-go-cli/internal/version"
)

func main() {
	cmd := newCmdBrev()
	if err := cmd.Execute(); err != nil {
		if _, ok := err.(*brev_errors.SuppressedError); ok {
			// error suppressed
		} else if brevError, ok := err.(brev_errors.BrevError); ok {
			fmt.Fprintln(os.Stderr, "Error: "+brevError.Error())
			fmt.Fprintln(os.Stderr, "\n"+brevError.Directive())
		} else {
			fmt.Fprint(os.Stderr, "Error: "+err.Error())
		}
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
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	brevCommand.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	brevCommand.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.Println(err)
		cmd.Println() // extra newline
		cmd.Println(cmd.UsageString())
		return &brev_errors.SuppressedError{}
	})

	createCmdTree(brevCommand, cmdContext)
	return brevCommand
}

func createCmdTree(brevCommand *cobra.Command, context *cmdcontext.Context) {
	brevCommand.AddCommand(endpoint.NewCmdEndpoint(context))
	brevCommand.AddCommand(auth.NewCmdLogin(context))
	brevCommand.AddCommand(package_project.NewCmdPackage(context))
	// brevCommand.AddCommand(project.NewCmdProject(context))
	brevCommand.AddCommand(version.NewCmdVersion(context))
	brevCommand.AddCommand(initialize.NewCmdInit(context))
	brevCommand.AddCommand(env.NewCmdEnv(context))
	brevCommand.AddCommand(status.NewCmdStatus(context))
	brevCommand.AddCommand(sync.NewCmdPull(context))
	brevCommand.AddCommand(sync.NewCmdPush(context))
	brevCommand.AddCommand(&completionCmd)
}
