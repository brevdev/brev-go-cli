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
package project

import (
	"github.com/spf13/cobra"

	"github.com/brevdev/brev-go-cli/internal/terminal"
)

func NewCmdProject(t *terminal.Terminal) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return logic(t)
		},
	}

	cmd.AddCommand(newCmdInit(t))
	cmd.AddCommand(newCmdList(t))
	cmd.AddCommand(newCmdLog(t))
	cmd.AddCommand(newCmdPull(t))
	cmd.AddCommand(newCmdPush(t))
	cmd.AddCommand(newCmdRemove(t))
	cmd.AddCommand(newCmdStatus(t))

	return cmd
}

func newCmdInit(t *terminal.Terminal) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return logic(t)
		},
	}

	return cmd
}

func newCmdList(t *terminal.Terminal) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return logic(t)
		},
	}

	return cmd
}

func newCmdLog(t *terminal.Terminal) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return logic(t)
		},
	}

	return cmd
}

func newCmdPull(t *terminal.Terminal) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return logic(t)
		},
	}

	return cmd
}

func newCmdPush(t *terminal.Terminal) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return logic(t)
		},
	}

	return cmd
}

func newCmdRemove(t *terminal.Terminal) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return logic(t)
		},
	}

	return cmd
}

func newCmdStatus(t *terminal.Terminal) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return logic(t)
		},
	}

	return cmd
}
