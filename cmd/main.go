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
	brevCommand.AddCommand(&completionCmd)
}
