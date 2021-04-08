package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/brevdev/brev-go-cli/internal/auth"
	"github.com/brevdev/brev-go-cli/internal/brev_errors"
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

func Test(c *cobra.Command) error {
	fmt.Println("heeyy")
	return nil
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
				t.Vprint(v)
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

	// 	help1 := `Usage:{{if .Runnable}}
	// 	{{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
	// 	{{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

	//   Aliases:
	// 	{{.NameAndAliases}}{{end}}{{if .HasExample}}

	//   Examples:
	//   {{.Example}}{{end}}{{if .HasAvailableSubCommands}}

	//   Flags:
	//   {{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

	//   Global Flags:
	//   {{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

	//   Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
	//   `
	// 	helpFinal := `
	//   Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
	//   `

	HouseKeeping := "\nHousekeeping Commands"
	Project := "\nProject Commands"
	Environment := "\nEnvironment Commands"
	Code := "\nCode Commands"

	for _, v := range brevCommand.Commands() {
		if v.Name() == "login" || v.Name() == "completion" {
			HouseKeeping += "\n\t" + v.Name() + "\t\t" + " " + v.Short
		} else if v.Name() == "init" || v.Name() == "clone" || v.Name() == "status" {
			Project += "\n\t" + v.Name() + "\t\t" + " " + v.Short
		} else if v.Name() == "env" || v.Name() == "packages" {
			Environment += "\n\t" + v.Name() + "\t\t" + " " + v.Short
		} else if v.Name() == "diff" || v.Name() == "endpoint" || v.Name() == "push" || v.Name() == "pull" {
			Code += "\n\t" + v.Name() + "\t\t" + " " + v.Short
		}
	}

	brevCommand.SetUsageTemplate(HouseKeeping + Project + Environment + Code)

	// 	brevCommand.SetUsageTemplate(`
	// {{t.Green('Usage:')}}
	// 	brev [flags]
	// 	brev [command]

	// Flags:
	// -h, --help      help for brev
	// -v, --verbose   Verbose output
	// --version   Print version output

	// HouseKeeping Commands:
	// 	login
	// 	completion

	// Project Commands:
	// 	init
	// 	clone
	// 	status

	// Environment Commands:
	// 	env
	// 	packages

	// Code Commands:
	// 	diff
	// 	endpoint
	// 	push
	// 	pull

	// `)

	// 	brevCommand.SetUsageTemplate(`Usage:{{if .Runnable}}
	// 	{{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
	// 	{{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

	//   Aliases:
	// 	{{.NameAndAliases}}{{end}}{{if .HasExample}}

	//   Examples:
	//   {{.Example}}{{end}}{{if .HasAvailableSubCommands}}

	//   Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
	// 	{{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

	//   Flags:
	//   {{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

	//   Global Flags:
	//   {{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

	//   Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
	// 	{{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

	//   Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
	//   `)

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
