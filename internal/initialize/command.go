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
package initialize

import (
	"fmt"
	"os"

	"github.com/brevdev/brev-go-cli/internal/auth"
	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
	"github.com/brevdev/brev-go-cli/internal/files"
	"github.com/spf13/cobra"
)

func NewCmdInit(context *cmdcontext.Context) *cobra.Command {
	var project string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a Brev Project",
		Long: `Use this to initialize a Brev project. Ex:
		
		// To init new project in current directory
		brev init
	
		// To init existing project
		brev init <project_name>
		`,
		RunE: func(cmd *cobra.Command, args []string) error {
			token, _ := auth.GetToken()
			brevAgent := brev_api.Agent{
				Key: token,
			}
			projects, err := brevAgent.GetProjects()
			if err != nil {
				context.PrintErr("Failed to retrieve projects", err)
				return err
			}

			for _, v := range projects {
				if v.Name == project {
					err = initExistingProj(v, context)
					if err != nil {
						return err
					}
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&project, "project", "p", "", "Project Name")
	cmd.RegisterFlagCompletionFunc("project", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getProjectNames(), cobra.ShellCompDirectiveNoSpace
	})

	return cmd
}

func getProjectNames() []string {

	// Get Projects
	token, _ := auth.GetToken()
	brevAgent := brev_api.Agent{
		Key: token,
	}
	rawProjects, _ := brevAgent.GetProjects()
	var projNames []string

	// Filter list for just project names
	for _, v := range rawProjects {
		projNames = append(projNames, v.Name)
	}

	// Return for shell completion
	return projNames
}

func initExistingProj(project brev_api.Project, context *cmdcontext.Context) error {

	// Get endpoints for project
	token, err := auth.GetToken()
	if err != nil {
		context.PrintErr("Failed to retrieve auth token", err)
		return err
	}

	brevAgent := brev_api.Agent{
		Key: token,
	}
	allEndpoints, err := brevAgent.GetEndpoints()
	if err != nil {
		context.PrintErr("Failed to get endpoints", err)
		return err
	}

	var endpoints brev_api.Endpoints
	for _, v := range allEndpoints {
		if v.ProjectId == project.Id {
			endpoints.Endpoints = append(endpoints.Endpoints, v)
		}
	}

	// Init the new folder at pwd + project name
	cwd, err := os.Getwd()
	if err != nil {
		context.PrintErr("Failed to determine working directory", err)
		return err
	}

	path := fmt.Sprintf("%s/%s", cwd, project.Name)

	// Make project.json
	err = files.OverwriteJSON(path+"/"+files.GetBrevDirectory()+"/"+files.GetProjectsFile(), project)
	if err != nil {
		context.PrintErr("Failed to write project to local file", err)
		return err
	}

	// Make endpoints.json
	err = files.OverwriteJSON(path+"/"+files.GetBrevDirectory()+"/"+files.GetEndpointsFile(), endpoints.Endpoints)
	if err != nil {
		context.PrintErr("Failed to write endpoints to local file", err)
		return err
	}

	// Create a global file with project directories
	var currBrevDirectories []string
	err = files.ReadJSON(files.GetActiveProjectsPath(), &currBrevDirectories)
	if err != nil {
		context.PrintErr("Failed to read projects directory", err)
		return err
	}

	if !brev_api.StringInList(path, currBrevDirectories) {
		currBrevDirectories = append(currBrevDirectories, path)
		err = files.OverwriteJSON(files.GetActiveProjectsPath(), currBrevDirectories)
		if err != nil {
			context.PrintErr("Failed to write projects to project file", err)
			return err
		}
	}

	// TODO: copy shared code

	// Create endpoint files
	for _, v := range endpoints.Endpoints {
		err = files.OverwriteString(fmt.Sprintf("%s/%s.py", path, v.Name), v.Code)
		if err != nil {
			context.PrintErr("Failed to write code to local file", err)
			return err
		}
	}

	return nil
}
