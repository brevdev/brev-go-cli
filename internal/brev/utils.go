package brev

import (
	"os"

	"github.com/brevdev/brev-go-cli/internal/auth"
	"github.com/brevdev/brev-go-cli/internal/config"
	"github.com/brevdev/brev-go-cli/internal/files"
)

type BrevAgent struct {
	Key *auth.CotterOauthToken
}

func brevEndpoint(resource string) string {
	baseEndpoint := config.GetBrevAPIEndpoint()
	return baseEndpoint + "/_api/" + resource
}

// Example usage
/*
	token, _ := auth.GetToken()
	brevAgent := brev.BrevAgent{
		Key: token,
	}

	endpointsResponse, _ := brevAgent.GetEndpoints()
	fmt.Println(endpointsResponse)

	projectsResponse, _ := brevAgent.GetProjects()
	fmt.Println(projectsResponse)

	modulesResponse, _ := brevAgent.GetModules()
	fmt.Println(modulesResponse)
*/

func GetActiveProject() BrevProject {
	cwd, _ := os.Getwd()
	path := cwd + "/.brev/projects.json"
	
	var project BrevProject
	_ = files.ReadJSON(path, &project)
	
	return project
}