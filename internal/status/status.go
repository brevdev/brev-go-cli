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
package status

import (
	"fmt"

	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
	"github.com/spf13/cobra"
)

func NewCmdStatus(context *cmdcontext.Context) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Get the latest project metadata",
		Long: `See high level on your project. Ex:

			brev status

		`, PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := cmdcontext.InvokeParentPersistentPreRun(cmd, args)
			if err != nil {
				return err
			}

			_, err = brev_api.CheckOutsideBrevErrorMessage(context)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("STATUS")
			return nil
		},
	}

	return cmd
}
