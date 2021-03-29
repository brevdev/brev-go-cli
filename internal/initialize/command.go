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
	"strings"

	"github.com/brevdev/brev-go-cli/internal/auth"
	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/brev_ctx"
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
	"github.com/brevdev/brev-go-cli/internal/files"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var green = color.New(color.FgGreen).SprintfFunc()
var yellow = color.New(color.FgYellow).SprintfFunc()
var red = color.New(color.FgRed).SprintfFunc()

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

			if project == "" {
				fmt.Fprint(context.VerboseOut, yellow("\nInitializing new project"))
			} else {
				fmt.Fprint(context.VerboseOut, yellow("\nInitializing project %s", project))
			}

			token, _ := auth.GetToken()
			brevAgent := brev_api.Agent{
				Key: token,
			}
			projects, err := brevAgent.GetProjects()
			if err != nil {
				context.PrintErr("Failed to retrieve projects", err)
				return err
			}

			if project == "" {
				err = initNewProject(context)
				if err != nil {
					return err
				}
			}

			for _, v := range projects {

				if v.Name == project {
					err = initExistingProj(v, context)
					if err != nil {
						return err
					}
					break // in case of error where multiple projects share name. We should prohibit this.
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

	cwd, err := os.Getwd()
	if err != nil {
		context.PrintErr(red("Failed to determine working directory"), err)
		return err
	}

	fmt.Fprint(context.VerboseOut, "\nCloning Brev project in "+yellow(cwd))
	fmt.Fprint(context.VerboseOut, green("\nCreating local files..."))

	// Get endpoints for project
	brevCtx, err := brev_ctx.New()
	if err != nil {
		return err
	}

	endpoints, err := brevCtx.Remote.GetEndpoints(&brev_ctx.GetEndpointsOptions{
		ProjectID: project.Id,
	})
	if err != nil {
		return err
	}

	// Init the new folder at pwd + project name
	path := fmt.Sprintf("%s/%s", cwd, project.Name)

	// Make project.json
	err = files.OverwriteJSON(path+"/"+files.GetBrevDirectory()+"/"+files.GetProjectsFile(), project)
	if err != nil {
		context.PrintErr(red("Failed to write project to local file"), err)
		return err
	}

	// Make endpoints.json
	err = files.OverwriteJSON(path+"/"+files.GetBrevDirectory()+"/"+files.GetEndpointsFile(), endpoints)
	if err != nil {
		context.PrintErr(red("Failed to write endpoints to local file"), err)
		return err
	}

	// Create a global file with project directories
	var currBrevDirectories []string
	err = files.ReadJSON(files.GetActiveProjectsPath(), &currBrevDirectories)
	if err != nil {
		context.PrintErr(red("Failed to read projects directory"), err)
		return err
	}

	if !brev_api.StringInList(path, currBrevDirectories) {
		currBrevDirectories = append(currBrevDirectories, path)
		err = files.OverwriteJSON(files.GetActiveProjectsPath(), currBrevDirectories)
		if err != nil {
			context.PrintErr(red("Failed to write projects to project file"), err)
			return err
		}
	}

	// TODO: copy shared code

	// Create endpoint files
	for _, v := range endpoints {
		err = files.OverwriteString(fmt.Sprintf("%s/%s.py", path, v.Name), v.Code)
		if err != nil {
			context.PrintErr("Failed to write code to local file", err)
			return err
		}
	}

	fmt.Fprint(context.VerboseOut, green("\n\nBrev project %s cloned.", project.Name))
	fmt.Fprint(context.VerboseOut, yellow("\ncd %s", project.Name))
	fmt.Fprint(context.VerboseOut, green(" and get started!"))
	fmt.Fprint(context.VerboseOut, green("\n\nHappy Hacking 🥞"))

	return nil
}

func initNewProject(context *cmdcontext.Context) error {

	// Get Project Name (parent folder-- behavior just like git init)
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	fmt.Fprint(context.VerboseOut, "\nCreating Brev project in "+yellow(cwd))

	dirs := strings.Split(cwd, "/")
	projName := dirs[len(dirs)-1]

	// Create new project
	token, _ := auth.GetToken()
	brevAgent := brev_api.Agent{
		Key: token,
	}
	projectResponse, _ := brevAgent.CreateProject(projName)
	project := projectResponse.Project

	fmt.Fprint(context.VerboseOut, green("\nCreating local files..."))

	// Make project.json
	err = files.OverwriteJSON(cwd+"/"+files.GetBrevDirectory()+"/"+files.GetProjectsFile(), project)
	if err != nil {
		context.PrintErr(red("Failed to write project to local file"), err)
		return err
	}

	// Make endpoints.json
	err = files.OverwriteJSON(cwd+"/"+files.GetBrevDirectory()+"/"+files.GetEndpointsFile(), []string{})
	if err != nil {
		context.PrintErr(red("Failed to write project to local file"), err)
		return err
	}

	// TODO: create shared code module

	// Add to path
	var currBrevDirectories []string
	err = files.ReadJSON(files.GetActiveProjectsPath(), &currBrevDirectories)
	if err != nil {
		context.PrintErr(red("Failed to read projects directory"), err)
		return err
	}

	if !brev_api.StringInList(cwd, currBrevDirectories) {
		currBrevDirectories = append(currBrevDirectories, cwd)
		err = files.OverwriteJSON(files.GetActiveProjectsPath(), currBrevDirectories)
		if err != nil {
			context.PrintErr(red("Failed to write projects to project file"), err)
			return err
		}
	}
	fmt.Fprint(context.VerboseOut, green("\n\nBrev project %s created and deployed.", projName))
	fmt.Fprint(context.VerboseOut, yellow("\ncd %s", projName))
	fmt.Fprint(context.VerboseOut, green(" and get started!"))
	fmt.Fprint(context.VerboseOut, green("\n\nHappy Hacking 🥞"))

	return nil

}
