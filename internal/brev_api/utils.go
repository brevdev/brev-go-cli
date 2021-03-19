package brev_api

import (
	"fmt"
	"os"
	"strings"

	"github.com/brevdev/brev-go-cli/internal/auth"
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
	"github.com/brevdev/brev-go-cli/internal/config"
	"github.com/brevdev/brev-go-cli/internal/files"
)

type Agent struct {
	Key *auth.CotterOauthToken
}

func brevEndpoint(resource string) string {
	baseEndpoint := config.GetBrevAPIEndpoint()
	return baseEndpoint + "/_api/" + resource
}

// Example usage
/*
	token, _ := auth.GetToken()
	brevAgent := brev_api.Agent{
		Key: token,
	}

	endpointsResponse, _ := brevAgent.GetEndpoints()
	fmt.Println(endpointsResponse)

	projectsResponse, _ := brevAgent.GetProjects()
	fmt.Println(projectsResponse)

	modulesResponse, _ := brevAgent.GetModules()
	fmt.Println(modulesResponse)
*/

func GetActiveProject() (*Project, error) {
	projectFilePath := files.GetProjectsPath()

	var project Project
	err := files.ReadJSON(projectFilePath, &project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func IsInProjectDirectory() (bool, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return false, err
	}

	var currBrevDirectories []string
	err = files.ReadJSON(files.GetActiveProjectsPath(), &currBrevDirectories)
	if err != nil {
		return false, err
	}

	for _, v := range currBrevDirectories {
		if strings.Contains(cwd, v) {
			return true, nil
		}
	}
	return false, nil
}

func CheckOutsideBrevErrorMessage(context *cmdcontext.Context) (bool, error) {
	isInProjectDirectory, err := IsInProjectDirectory()
	if err != nil {
		return false, nil
	}

	if isInProjectDirectory {
		return true, nil
	}

	var currBrevDirectories []string
	err = files.ReadJSON(files.GetActiveProjectsPath(), &currBrevDirectories)
	if err != nil {
		context.PrintErr("Failed to read projects from local directory", err)
		return false, err
	}

	// Exit with error message
	fmt.Fprintln(context.Out, "Endpoint commands only work in a Brev project.")
	if len(currBrevDirectories) == 0 {
		// If no directories, check if they have some remote.

		// Get Projects
		token, err := auth.GetToken()
		if err != nil {
			context.PrintErr("Failed to retrieve auth token", err)
			return false, err
		}
		brevAgent := Agent{
			Key: token,
		}
		rawProjects, err := brevAgent.GetProjects()
		if err != nil {
			context.PrintErr("Failed to get projects", err)
			return false, err
		}

		if len(rawProjects) == 0 {
			// Encourage them to create their first project
			fmt.Fprintln(context.Out, "You haven't made a brev project yet! Try running 'brev init'")

		} else {
			// Encourage them to pull one of their existing projects
			fmt.Fprintln(context.Out, "Set up one of your existing projects.")
			fmt.Fprintln(context.Out, "For example, run 'brev init "+rawProjects[0].Name+"'")
		}

	} else {
		// Print active brev projects
		fmt.Fprintln(context.Out, "Active Brev projects on your computer: ")
		for _, v := range currBrevDirectories {
			fmt.Fprintln(context.Out, "\t"+v)
		}
	}
	return false, nil
}

func StringInList(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
