package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/brevdev/brev-go-cli/internal/auth"
	"github.com/brevdev/brev-go-cli/internal/brev_errors"
	"github.com/brevdev/brev-go-cli/internal/endpoint"
	"github.com/brevdev/brev-go-cli/internal/env"
	"github.com/brevdev/brev-go-cli/internal/initialize"
	"github.com/brevdev/brev-go-cli/internal/logs"
	"github.com/brevdev/brev-go-cli/internal/package_project"
	"github.com/brevdev/brev-go-cli/internal/status"
	"github.com/brevdev/brev-go-cli/internal/sync"
	"github.com/brevdev/brev-go-cli/internal/terminal"
	"github.com/brevdev/brev-go-cli/internal/version"
)

func main() {
	t := &terminal.Terminal{}

	cmd := newCmdBrev(t)
	if err := cmd.Execute(); err != nil {
		if _, ok := err.(*brev_errors.SuppressedError); ok {
			// error suppressed
		} else {
			t.Errprint(err, "")
		}
		os.Exit(1)
	}
}

func newCmdBrev(t *terminal.Terminal) *cobra.Command {
	var verbose bool
	var printVersion bool

	brevCommand := &cobra.Command{
		Use: "brev",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			t.Init(verbose)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if printVersion {
				v, err := version.BuildVersionString(t)
				if err != nil {
					t.Errprint(err, "Failed to determine version")
					return err
				}
				t.Println(v)
				return nil
			} else {
				return cmd.Usage()
			}
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	brevCommand.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	brevCommand.PersistentFlags().BoolVar(&printVersion, "version", false, "Print version output")
	brevCommand.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.Println(err)
		cmd.Println() // extra newline
		cmd.Println(cmd.UsageString())
		return &brev_errors.SuppressedError{}
	})

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
	brevCommand.AddCommand((logs.NewCmdLogs(t)))
	brevCommand.AddCommand(&completionCmd)
}
