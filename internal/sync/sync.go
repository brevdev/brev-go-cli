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
package sync

import (
	"github.com/spf13/cobra"

	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
	"github.com/brevdev/brev-go-cli/internal/terminal"
)

func NewCmdPush(t *terminal.Terminal) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push",
		Short: "Push your local changes to remote",
		Long: `To push your local changes:

			brev push
		`, PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := cmdcontext.InvokeParentPersistentPreRun(cmd, args)
			if err != nil {
				return err
			}

			_, err = brev_api.CheckOutsideBrevErrorMessage(t)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return push(t)
		},
	}

	return cmd
}

func NewCmdPull(t *terminal.Terminal) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull latest changes from your server",
		Long: `To pull latest changes:

			brev pull
		`, PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := cmdcontext.InvokeParentPersistentPreRun(cmd, args)
			if err != nil {
				return err
			}

			_, err = brev_api.CheckOutsideBrevErrorMessage(t)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return pull(t)
		},
	}

	return cmd
}

func NewCmdDiff(t *terminal.Terminal) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "diff",
		Short: "See a diff of your local changes compared to what's deployed in the console",
		Long: `To see a diff of your local changes compared to what's deployed in the console,
			from an active brev project directory, run:

			brev diff
		`, PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := cmdcontext.InvokeParentPersistentPreRun(cmd, args)
			if err != nil {
				return err
			}

			_, err = brev_api.CheckOutsideBrevErrorMessage(t)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return diffCmd(t)
		},
	}

	return cmd
}
