package initialize

import (
	"errors"
	"fmt"
	"os"

	"github.com/brevdev/brev-go-cli/internal/auth"
	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/brev_ctx"
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
)

func initializeExistingProject(projectName string, context *cmdcontext.Context) error {
	project, endpoints, err := getRemoteProjectMatchingName(projectName)
	if err != nil {
		context.PrintErr(fmt.Sprintf("Failed to retrieve project with name '%s': %s", projectName, err), err)
		return err
	}

	// Set local
	err = brev_ctx.SetLocal(&brev_ctx.LocalContext{
		Project:   *project,
		Endpoints: endpoints.Endpoints,
	})
	if err != nil {
		return errors.New("Failed to update local brev project")
	}

	// Update global
	cwd, _ := os.Getwd()
	err = brev_ctx.SetGlobalProjectPath(cwd + "/" + projectName)
	if err != nil {
		return errors.New("failed to update global brev state")
	}

	return nil
}

func getRemoteProjectMatchingName(projectName string) (*brev_api.Project, *brev_api.Endpoints, error) {
	// Get remote projects
	token, _ := auth.GetToken()
	brevAgent := brev_api.Agent{
		Key: token,
	}
	remoteProjects, err := brevAgent.GetProjects()
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Failed to retrieve remote projects: %s", err))
	}

	// Filter
	var project *brev_api.Project
	for _, remoteProject := range remoteProjects {
		if remoteProject.Name == projectName {
			project = &remoteProject
		}
	}
	if project == nil {
		return nil, nil, errors.New(fmt.Sprintf("Failed to find project matching name '%s'", projectName))
	}

	// Get endpoints for project
	remoteEndpoints, err := brevAgent.GetEndpoints()
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Failed to retrieve remote endpoints: %s", err))
	}

	// Filter
	var endpoints brev_api.Endpoints
	for _, remoteEndpoint := range remoteEndpoints.Endpoints {
		if remoteEndpoint.ProjectId == project.Id {
			endpoints.Endpoints = append(endpoints.Endpoints, remoteEndpoint)
		}
	}

	return project, &endpoints, nil
}
