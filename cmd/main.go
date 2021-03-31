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
	"github.com/brevdev/brev-go-cli/internal/endpoint"
	"github.com/brevdev/brev-go-cli/internal/env"
	"github.com/brevdev/brev-go-cli/internal/initialize"
	"github.com/brevdev/brev-go-cli/internal/package_project"
	"github.com/brevdev/brev-go-cli/internal/status"
	"github.com/brevdev/brev-go-cli/internal/sync"
	"github.com/brevdev/brev-go-cli/internal/terminal"
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
	var print_version bool
	t := &terminal.Terminal{}

	brevCommand := &cobra.Command{
		Use: "brev",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			fmt.Print("\n")
			t.Init(verbose)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if print_version {
				v, err := version.BuildVersionString(t)
				if err != nil {
					t.Errprint(err, "Failed to determine version")
					return err
				}
				t.Vprint(v)
				return nil
			} else {
				return cmd.Usage()
			}
		},
	}

	brevCommand.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	brevCommand.PersistentFlags().BoolVar(&print_version, "version", false, "Print version output")

	createCmdTree(brevCommand, t)
	return brevCommand
}

func createCmdTree(brevCommand *cobra.Command, t *terminal.Terminal) {
	brevCommand.AddCommand(endpoint.NewCmdEndpoint(t))
	brevCommand.AddCommand(auth.NewCmdLogin(t))
	brevCommand.AddCommand(package_project.NewCmdPackage(t))
	// brevCommand.AddCommand(project.NewCmdProject(context))
	brevCommand.AddCommand(initialize.NewCmdInit(t))
	brevCommand.AddCommand(env.NewCmdEnv(t))
	brevCommand.AddCommand(status.NewCmdStatus(t))
	brevCommand.AddCommand(sync.NewCmdPull(t))
	brevCommand.AddCommand(sync.NewCmdPush(t))
	brevCommand.AddCommand(sync.NewCmdDiff(t))
	brevCommand.AddCommand(&completionCmd)
}
