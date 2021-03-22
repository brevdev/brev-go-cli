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
	brevCtx, err := brev_ctx.New()
	if err != nil {
		return err
	}

	if err = brevCtx.Local.SetProject(*project); err != nil {
		return err
	}
	if err = brevCtx.Local.SetEndpoints(endpoints.Endpoints); err != nil {
		return err
	}

	// Update global
	cwd, _ := os.Getwd()
	if err = brevCtx.Global.SetProjectPath(cwd + "/" + projectName); err != nil {
		return err
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
	for _, remoteEndpoint := range remoteEndpoints {
		if remoteEndpoint.ProjectId == project.Id {
			endpoints.Endpoints = append(endpoints.Endpoints, remoteEndpoint)
		}
	}

	return project, &endpoints, nil
}
