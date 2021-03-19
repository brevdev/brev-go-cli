package brev_ctx

import (
	"errors"
	"fmt"
	"os"

	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/files"
)

type GlobalContext struct {
	ProjectPaths []string
}

type LocalContext struct {
	Project   brev_api.Project
	Endpoints []brev_api.Endpoint
}

const (
	localBrevDirectory  = ".brev"
	globalBrevDirectory = ".brev"

	localProjectsFile        = "projects.json"
	localEndpointsFile       = "endpoints.json"
	globalActiveProjectsFile = "active_projects.json"
)

// GetLocal returns the global brev context.
func GetGlobal() (*GlobalContext, error) {
	globalActiveProjectsFileExists, err := files.Exists(getGlobalActiveProjectsPath())
	if err != nil {
		return nil, errors.New(fmt.Sprintf(
			"Failed to read from %s: %s", getGlobalActiveProjectsPath(), err,
		))
	}
	if !globalActiveProjectsFileExists {
		return &GlobalContext{}, nil
	}

	var brevProjectPaths []string
	err = files.ReadJSON(getGlobalActiveProjectsPath(), &brevProjectPaths)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(
			"Failed to read from %s: %s", getGlobalActiveProjectsPath(), err,
		))
	}
	return &GlobalContext{
		ProjectPaths: brevProjectPaths,
	}, nil
}

// SetGlobal overwrites the global brev context.
func SetGlobal(context *GlobalContext) error {
	err := files.OverwriteJSON(getGlobalActiveProjectsPath(), context.ProjectPaths)
	if err != nil {
		return errors.New(fmt.Sprintf(
			"Failed to write to %s: %s", getGlobalActiveProjectsPath(), err,
		))
	}
	return nil
}

// SetGlobalProjectPath adds (or replaces) the given project path in the blobal brev context.
func SetGlobalProjectPath(path string) error {
	global, err := GetGlobal()
	if err != nil {
		return err
	}
	var exists bool
	for _, globalProjectPath := range global.ProjectPaths {
		if path == globalProjectPath {
			exists = true
		}
	}
	if !exists {
		global.ProjectPaths = append(global.ProjectPaths, path)
	}
	return SetGlobal(global)
}

// GetLocal returns the local brev context. If the brev context is nil, the current
// directory is not a bre project.
func GetLocal() (*LocalContext, error) {
	localProjectFileExists, err := files.Exists(getLocalProjectPath())
	if err != nil {
		return nil, errors.New(fmt.Sprintf(
			"Failed to read from %s: %s", getLocalProjectPath(), err,
		))
	}
	localEndpointsFileExists, err := files.Exists(getLocalEndpointsPath())
	if err != nil {
		return nil, errors.New(fmt.Sprintf(
			"Failed to read from %s: %s", getLocalEndpointsPath(), err,
		))
	}

	if !localProjectFileExists && localEndpointsFileExists {
		// endpoints.json exists, but projects.json does not
		return nil, errors.New(fmt.Sprintf(
			"Local project is corrupted: Failed to read from %s: %s", getLocalProjectPath(), err,
		))
	} else if !localEndpointsFileExists && localProjectFileExists {
		// projects.json exists, but endpoints.json does not
		return nil, errors.New(fmt.Sprintf(
			"Local project is corrupted: Failed to read from %s: %s", getLocalEndpointsPath(), err,
		))
	} else if !localEndpointsFileExists && !localProjectFileExists {
		// acceptable case -- no files exist
		return nil, nil
	}

	var brevProject brev_api.Project
	err = files.ReadJSON(getLocalProjectPath(), &brevProject)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(
			"Failed to read from %s: %s", getLocalProjectPath(), err,
		))
	}

	var brevEndpoints []brev_api.Endpoint
	err = files.ReadJSON(getLocalEndpointsPath(), &brevEndpoints)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(
			"Failed to read from %s: %s", getLocalEndpointsPath(), err,
		))
	}

	return &LocalContext{
		Project:   brevProject,
		Endpoints: brevEndpoints,
	}, nil
}

// SetLocal overwrites the local brev context.
func SetLocal(context *LocalContext) error {
	err := files.OverwriteJSON(getLocalProjectPath(), context.Project)
	if err != nil {
		return errors.New(fmt.Sprintf(
			"Failed to write to %s: %s", getLocalProjectPath(), err,
		))
	}
	err = files.OverwriteJSON(getLocalEndpointsPath(), context.Endpoints)
	if err != nil {
		return errors.New(fmt.Sprintf(
			"Failed to write to %s: %s", getLocalEndpointsPath(), err,
		))
	}
	return nil
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
