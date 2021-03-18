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
	"github.com/brevdev/brev-go-cli/internal/brev"
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
	"github.com/brevdev/brev-go-cli/internal/files"
	"github.com/spf13/cobra"
)

func NewCmdInit(context *cmdcontext.Context) *cobra.Command {
	var project string;

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a Brev Project",
		Long: `Use this to initialize a Brev project. Ex:
		
		// To init new project in current directory
		brev init
	
		// To init existing project
		brev init <project_name>
		`,
		Run: func(cmd *cobra.Command, args []string) {
			token, _ := auth.GetToken()
			brevAgent := brev.BrevAgent{
				Key: token,
			}
			projects, _ := brevAgent.GetProjects()

			for _, v := range projects {
				if (v.Name==project) {
					init_existing_proj(v)
				}
			}			
		},
	
	}
	cmd.Flags().StringVarP(&project, "project", "p", "", "Project Name")
	cmd.RegisterFlagCompletionFunc("project", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return get_project_names(), cobra.ShellCompDirectiveNoSpace})
	

	return cmd
}

func get_project_names() []string {

	// Get Projects
	token, _ := auth.GetToken()
	brevAgent := brev.BrevAgent{
		Key: token,
	}
	raw_projects, _ := brevAgent.GetProjects()
	var projNames []string

	// Filter list for just project names
	for _, v := range raw_projects {
		projNames = append(projNames, v.Name)
	}
	
	// Return for shell completion
	return projNames
}


func init_existing_proj(project brev.BrevProject) {

	// TODO: create a global file with project directories

	// Init the new folder at pwd + project name
	cwd, _ := os.Getwd()
	path := fmt.Sprintf("%s/%s/.brev", cwd, project.Name)

	// Make project.json 
	files.OverwriteJSON(path+"/projects.json", project)


	// Get endpoints for project
	token, _ := auth.GetToken()
	brevAgent := brev.BrevAgent{
		Key: token,
	}
	all_endpoints, _ := brevAgent.GetEndpoints()
	var endpoints brev.BrevEndpoints
	for _, v := range all_endpoints.Endpoints {
		if (v.ProjectId==project.Id) {
			endpoints.Endpoints = append(endpoints.Endpoints, v)
		}
	}	

	// Make endpoints.json
	files.OverwriteJSON(path+"/endpoints.json", endpoints.Endpoints)

	// TODO: copy shared code
	// TODO: copy endpoint files
	// TODO: copy variables as file ... ? should we do this?
	
}
