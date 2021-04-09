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
	t := terminal.New()

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
			t.SetVerbose(verbose)
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

	cobra.AddTemplateFunc("hasHousekeepingCommands", hasHousekeepingCommands)
	cobra.AddTemplateFunc("isHousekeepingCommand", isHousekeepingCommand)
	cobra.AddTemplateFunc("housekeepingCommands", housekeepingCommands)
	cobra.AddTemplateFunc("hasProjectCommands", hasProjectCommands)
	cobra.AddTemplateFunc("isProjectCommand", isProjectCommand)
	cobra.AddTemplateFunc("projectCommands", projectCommands)
	cobra.AddTemplateFunc("hasEnvironmentCommands", hasEnvironmentCommands)
	cobra.AddTemplateFunc("isEnvironmentCommand", isEnvironmentCommand)
	cobra.AddTemplateFunc("environmentCommands", environmentCommands)
	cobra.AddTemplateFunc("hasCodeCommands", hasCodeCommands)
	cobra.AddTemplateFunc("isCodeCommand", isCodeCommand)
	cobra.AddTemplateFunc("codeCommands", codeCommands)
	brevCommand.SetUsageTemplate(usageTemplate)

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
	brevCommand.AddCommand(initialize.NewCmdClone(t))
	brevCommand.AddCommand(initialize.NewCmdInit(t))
	brevCommand.AddCommand(env.NewCmdEnv(t))
	brevCommand.AddCommand(status.NewCmdStatus(t))
	brevCommand.AddCommand(sync.NewCmdPull(t))
	brevCommand.AddCommand(sync.NewCmdPush(t))
	brevCommand.AddCommand(sync.NewCmdDiff(t))
	brevCommand.AddCommand((logs.NewCmdLogs(t)))
	brevCommand.AddCommand(&completionCmd)
}

func hasHousekeepingCommands(cmd *cobra.Command) bool {
	return len(housekeepingCommands(cmd)) > 0
}

func hasProjectCommands(cmd *cobra.Command) bool {
	return len(projectCommands(cmd)) > 0
}

func hasEnvironmentCommands(cmd *cobra.Command) bool {
	return len(environmentCommands(cmd)) > 0
}

func hasCodeCommands(cmd *cobra.Command) bool {
	return len(codeCommands(cmd)) > 0
}

func housekeepingCommands(cmd *cobra.Command) []*cobra.Command {
	cmds := []*cobra.Command{}
	for _, sub := range cmd.Commands() {
		if isHousekeepingCommand(sub) {
			cmds = append(cmds, sub)
		}
	}
	return cmds
}

func projectCommands(cmd *cobra.Command) []*cobra.Command {
	cmds := []*cobra.Command{}
	for _, sub := range cmd.Commands() {
		if isProjectCommand(sub) {
			cmds = append(cmds, sub)
		}
	}
	return cmds
}

func environmentCommands(cmd *cobra.Command) []*cobra.Command {
	cmds := []*cobra.Command{}
	for _, sub := range cmd.Commands() {
		if isEnvironmentCommand(sub) {
			cmds = append(cmds, sub)
		}
	}
	return cmds
}

func codeCommands(cmd *cobra.Command) []*cobra.Command {
	cmds := []*cobra.Command{}
	for _, sub := range cmd.Commands() {
		if isCodeCommand(sub) {
			cmds = append(cmds, sub)
		}
	}
	return cmds
}

func isHousekeepingCommand(cmd *cobra.Command) bool {
	if _, ok := cmd.Annotations["housekeeping"]; ok {
		return true
	} else {
		return false
	}
}

func isProjectCommand(cmd *cobra.Command) bool {
	if _, ok := cmd.Annotations["project"]; ok {
		return true
	} else {
		return false
	}
}

func isEnvironmentCommand(cmd *cobra.Command) bool {
	if _, ok := cmd.Annotations["environment"]; ok {
		return true
	} else {
		return false
	}
}

func isCodeCommand(cmd *cobra.Command) bool {
	if _, ok := cmd.Annotations["code"]; ok {
		return true
	} else {
		return false
	}
}

var usageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

{{- if hasHousekeepingCommands . }}

Housekeeping Commands:
{{- range housekeepingCommands . }}
  {{rpad .Name .NamePadding }} {{.Short}}
{{- end}}{{- end}}{{- if hasProjectCommands . }}

Environment Commands:
{{- range environmentCommands . }}
  {{rpad .Name .NamePadding }} {{.Short}}
{{- end}}{{- end}}{{- if hasCodeCommands . }}

Project Commands:
{{- range projectCommands . }}
  {{rpad .Name .NamePadding }} {{.Short}}
{{- end}}{{- end}}{{- if hasEnvironmentCommands . }}

Code Commands:
{{- range codeCommands . }}
  {{rpad .Name .NamePadding }} {{.Short}}
{{- end}}{{- end}}{{- end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
