package brev_ctx

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/brevdev/brev-go-cli/internal/auth"
	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/files"
)

const (
	localBrevDirectory  = ".brev"
	globalBrevDirectory = ".brev"

	localProjectsFile        = "projects.json"
	localEndpointsFile       = "endpoints.json"
	globalActiveProjectsFile = "active_projects.json"
)

// BrevContext bundles all contexts
type BrevContext struct {
	Local  *LocalContext
	Global *GlobalContext
	Remote *RemoteContext
}

// GlobalContext encapsulates the Brev state of the current machine (e.g. laptop)
type GlobalContext struct{}

// LocalContext encapsulates the Brev state of the project corresponding to the current working directory
type LocalContext struct{}

// RemoteContext encapsulates the remote Brev state corresponding to the authorized user
type RemoteContext struct {
	agent *brev_api.Agent
}

type GetEndpointsOptions struct {
	ID        string
	Name      string
	ProjectID string
}

type GetProjectsOptions struct {
	ID   string
	Name string
}

type GetVariablesOptions struct {
	Name string
}

type GetPackagesOptions struct {
	Name string
}

// New instantiates a new instance of a BrevContext
func New() (*BrevContext, error) {
	local, err := NewLocal()
	if err != nil {
		return nil, fmt.Errorf("could not instantiate local context: %s", err)
	}
	global, err := NewGlobal()
	if err != nil {
		return nil, fmt.Errorf("could not instantiate global context: %s", err)
	}
	remote, err := NewRemote()
	if err != nil {
		return nil, fmt.Errorf("could not instantiate remote context: %s", err)
	}
	return &BrevContext{
		Local:  local,
		Global: global,
		Remote: remote,
	}, nil
}

// NewGlobal creates a new instance of a GlobalContext
func NewGlobal() (*GlobalContext, error) {
	return &GlobalContext{}, nil
}

// GetProjectPaths returns the filepaths of all projects known to the current system
func (c *GlobalContext) GetProjectPaths() ([]string, error) {
	globalActiveProjectsFileExists, err := files.Exists(getGlobalActiveProjectsPath())
	if err != nil {
		return nil, fmt.Errorf("failed to read from %s: %s", getGlobalActiveProjectsPath(), err)
	}
	if !globalActiveProjectsFileExists {
		return nil, nil
	}

	var paths []string
	err = files.ReadJSON(getGlobalActiveProjectsPath(), &paths)
	if err != nil {
		return nil, fmt.Errorf("failed to read from %s: %s", getGlobalActiveProjectsPath(), err)
	}
	return paths, nil
}

// SetProjectPath sets the given path into the global list of project filepaths
func (c *GlobalContext) SetProjectPath(path string) error {
	paths, err := c.GetProjectPaths()
	if err != nil {
		return fmt.Errorf("failed to get project paths: %s", err)
	}

	for _, savedPath := range paths {
		if path == savedPath {
			// already exists -- return early
			return nil
		}
	}
	paths = append(paths, path)

	err = files.OverwriteJSON(getGlobalActiveProjectsPath(), paths)
	if err != nil {
		return fmt.Errorf("failed to write to %s: %s", getGlobalActiveProjectsPath(), err)
	}
	return nil
}

// NewLocal creates a new instance of a LocalContext
func NewLocal() (*LocalContext, error) {
	return &LocalContext{}, nil
}

// GetEndpoints returns a list of endpoints for the current local brev project. An optional
// GetEndpointsOptions struct may be provided to filter the results.
//
// Example usage:
//   local, _ := NewLocal()
//
//   // no filtering
//   endpoints, err := local.GetEndpoints(nil)
//
//   // filter by endpoint name
//   endpointsForName, _ := remote.GetEndpoints(&GetProjectsOptions{
//       Name: "foobarbaz",
//   })
func (c *LocalContext) GetEndpoints(options *GetEndpointsOptions) ([]brev_api.Endpoint, error) {
	localEndpointsFileExists, err := files.Exists(getLocalEndpointsPath())
	if err != nil {
		return nil, errors.New(fmt.Sprintf(
			"Failed to read from %s: %s", getLocalEndpointsPath(), err,
		))
	}
	if !localEndpointsFileExists {
		return nil, nil
	}

	var endpoints []brev_api.Endpoint
	err = files.ReadJSON(getLocalEndpointsPath(), &endpoints)
	if err != nil {
		return nil, fmt.Errorf("failed to read from %s: %s", getLocalEndpointsPath(), err)
	}

	if options == nil {
		return endpoints, nil
	}

	var filteredEndpoints []brev_api.Endpoint
	for _, endpoint := range endpoints {
		if options.ID != "" && endpoint.Id == options.ID {
			filteredEndpoints = append(filteredEndpoints, endpoint)
			break
		}
		if options.Name != "" && endpoint.Name == options.Name {
			filteredEndpoints = append(filteredEndpoints, endpoint)
		}
		if options.ProjectID != "" && endpoint.ProjectId == options.ProjectID {
			filteredEndpoints = append(filteredEndpoints, endpoint)
		}
	}
	return filteredEndpoints, nil
}

// GetProject retrieves the project associated with the current working directory
func (c *LocalContext) GetProject() (*brev_api.Project, error) {
	localProjectFileExists, err := files.Exists(getLocalProjectPath())
	if err != nil {
		return nil, fmt.Errorf("failed to read from %s: %s", getLocalProjectPath(), err)
	}
	if !localProjectFileExists {
		return nil, nil
	}
	var project brev_api.Project
	err = files.ReadJSON(getLocalProjectPath(), &project)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(
			"Failed to read from %s: %s", getLocalProjectPath(), err,
		))
	}
	return &project, nil
}

// SetProject stores the state of the given project in the context of the current working directory
func (c *LocalContext) SetProject(project brev_api.Project) error {
	err := files.OverwriteJSON(getLocalProjectPath(), project)
	if err != nil {
		return fmt.Errorf("failed to write to %s: %s", getLocalProjectPath(), err)
	}
	return nil
}

// SetEndpoints stores the state of the given endpoints in the context of the current working directory
func (c *LocalContext) SetEndpoints(endpoints []brev_api.Endpoint) error {
	err := files.OverwriteJSON(getLocalEndpointsPath(), endpoints)
	if err != nil {
		return fmt.Errorf("failed to write to %s: %s", getLocalEndpointsPath(), err)
	}
	return nil
}

func (c *LocalContext) SetEndpoint(endpoint brev_api.Endpoint) error {
	endpoints, err := c.GetEndpoints(nil)
	if err != nil {
		return nil
	}

	// if endpoint is new, save
	var exists bool
	for _, savedEndpoint := range endpoints {
		if reflect.DeepEqual(endpoint, savedEndpoint) {
			exists = true
		}
	}
	if !exists {
		endpoints = append(endpoints, endpoint)
		err = files.OverwriteJSON(getLocalEndpointsPath(), endpoints)
		if err != nil {
			return fmt.Errorf("failed to write to %s: %s", getLocalEndpointsPath(), err)
		}
	}

	return nil
}

// NewRemote returns a new instance of a RemoteContext, with an initialized auth token.
// Further calls to NewRemote will re-authenticate and store new auth tokens.
func NewRemote() (*RemoteContext, error) {
	token, err := auth.GetToken()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve auth token: %s", err)
	}

	return &RemoteContext{
		agent: &brev_api.Agent{Key: token},
	}, nil
}

// GetProjects retrieves remote projects for the context user. An optional GetProjectsOptions
// struct may be provided to filter the results.
//
// Example usage:
//   remote, _ := NewRemote()
//
//   // no filtering
//   projects, err := remote.GetProjects(nil)
//
//   // filter by project name
//   projectsForName, _ := remote.GetProjects(&GetProjectsOptions{
//       Name: "foobarbaz",
//   })
//
//   // filter by project ID
//   projectsForID, _ := remote.GetProjects(&GetProjectsOptions{
//       ID: "abc123def456",
//   })
func (c *RemoteContext) GetProjects(options *GetProjectsOptions) ([]brev_api.Project, error) {
	projects, err := c.agent.GetProjects()
	if err != nil {
		return nil, fmt.Errorf("failed to get endpoints: %s", err)
	}

	if options == nil {
		return projects, nil
	}

	var filteredProjects []brev_api.Project
	for _, project := range projects {
		if options.ID != "" && project.Id == options.ID {
			filteredProjects = append(filteredProjects, project)
			break
		}
		if options.Name != "" && project.Name == options.Name {
			filteredProjects = append(filteredProjects, project)
		}
	}
	return filteredProjects, nil
}

// GetEndpoints retrieves remote endpoints for the context user. An optional GetEndpointsOptions
// struct may be provided to filter the results.
//
// Example usage:
//   remote, _ := NewRemote()
//
//   // no filtering
//   endpoints, err := remote.GetEndpoints(nil)
//
//   // filter by endpoint name
//   endpointsForName, _ := remote.GetEndpoints(&GetEndpointsOptions{
//       Name: "foobarbaz",
//   })
//
//   // filter by project ID
//   endpointsForProject, _ := remote.GetEndpoints(&GetEndpointsOptions{
//       ProjectID: "abc123def456",
//   })
func (c *RemoteContext) GetEndpoints(options *GetEndpointsOptions) ([]brev_api.Endpoint, error) {
	endpoints, err := c.agent.GetEndpoints()
	if err != nil {
		return nil, fmt.Errorf("failed to get endpoints: %s", err)
	}

	if options == nil {
		return endpoints, nil
	}

	var filteredEndpoints []brev_api.Endpoint
	for _, endpoint := range endpoints {
		if options.ID != "" && endpoint.Id == options.ID {
			filteredEndpoints = append(filteredEndpoints, endpoint)
			break
		}
		if options.Name != "" && endpoint.Name == options.Name {
			filteredEndpoints = append(filteredEndpoints, endpoint)
		}
		if options.ProjectID != "" && endpoint.ProjectId == options.ProjectID {
			filteredEndpoints = append(filteredEndpoints, endpoint)
		}
	}
	return filteredEndpoints, nil
}

// SetEndpoint updates the remote endpoint with the given ID with the state of the provided
// endpoint struct.
func (c *RemoteContext) SetEndpoint(endpoint brev_api.Endpoint) (*brev_api.Endpoint, error) {
	if endpoint.Id == "" {
		response, err := c.agent.CreateEndpoint(endpoint.Name, endpoint.ProjectId)
		if err != nil {
			return nil, fmt.Errorf("failed to create endpoint: %s", err)
		}
		return &response.Endpoint, nil
	} else {
		response, err := c.agent.UpdateEndpoint(endpoint.Id, brev_api.RequestUpdateEndpoint{
			Name:    endpoint.Name,
			Methods: endpoint.Methods,
			Code:    endpoint.Code,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to update endpoint: %s", err)
		}
		return &response.Endpoint, nil
	}
}

// DeleteEndpoint removes the remote endpoint with the given ID.
func (c *RemoteContext) DeleteEndpoint(endpoint brev_api.Endpoint) error {
	_, err := c.agent.RemoveEndpoint(endpoint.Id)
	if err != nil {
		return fmt.Errorf("failed to delete endpoint: %s", err)
	}
	return nil
}

// GetVariables retrieves remote variables for the given project. An optional GetVariablesOptions
// struct may be provided to filter the results.
//
// Example usage:
//   remote, _ := NewRemote()
//
//   // no filtering
//   endpoints, err := remote.GetVariables(nil)
//
//   // filter by variable name
//   variablesForProject, _ := remote.GetVariables(&GetVariablesOptions{
//       Name: "foobarbaz",
//   })
func (c *RemoteContext) GetVariables(project brev_api.Project, options *GetVariablesOptions) ([]brev_api.ProjectVariable, error) {
	variables, err := c.agent.GetVariables(project.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get project variables: %s", err)
	}

	if options == nil {
		return variables, nil
	}

	var filteredVariables []brev_api.ProjectVariable
	for _, variable := range variables {
		if options.Name != "" && variable.Name == options.Name {
			filteredVariables = append(filteredVariables, variable)
			break
		}
	}
	return filteredVariables, nil
}

// SetVariable sets the given name/value pair for the given project.
func (c *RemoteContext) SetVariable(project brev_api.Project, name string, value string) (*brev_api.ProjectVariable, error) {
	response, err := c.agent.AddVariable(project.Id, name, value)
	if err != nil {
		return nil, fmt.Errorf("failed to set project variable: %s", err)
	}
	return &response.Variable, nil
}

func (c *RemoteContext) GetPackages(project brev_api.Project, options *GetPackagesOptions) ([]brev_api.ProjectPackage, error) {
	packages, err := c.agent.GetPackages(project.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve project packages: %s", err)
	}

	if options == nil {
		return packages, nil
	}
	var filteredPackages []brev_api.ProjectPackage
	for _, projectPackage := range packages {
		if options.Name != "" && projectPackage.Name == options.Name {
			filteredPackages = append(filteredPackages, projectPackage)
			break
		}
	}
	return filteredPackages, nil
}

func (c *RemoteContext) SetPackage(project brev_api.Project, name string) (*brev_api.ProjectPackage, error) {
	response, err := c.agent.AddPackage(project.Id, name)
	if err != nil {
		return nil, fmt.Errorf("failed to add project package: %s", err)
	}
	return &response.Package, nil
}

func getGlobalActiveProjectsPath() string {
	homeDir := files.GetHomeDir()
	return fmt.Sprintf("%s/%s/%s", homeDir, globalBrevDirectory, globalActiveProjectsFile)
}

func getLocalProjectPath() string {
	cwd, _ := os.Getwd()
	return fmt.Sprintf("%s/%s/%s", cwd, localBrevDirectory, localProjectsFile)
}

func getLocalEndpointsPath() string {
	cwd, _ := os.Getwd()
	return fmt.Sprintf("%s/%s/%s", cwd, localBrevDirectory, localEndpointsFile)
}
