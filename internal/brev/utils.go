package brev

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"

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

func GetRootDir() string {
	usr, err := user.Current()
    if err != nil {
        log.Fatal( err )
    }
    return usr.HomeDir
}


func IsInProjectDirectory() bool {
	cwd,_ := os.Getwd()
	rootDir := GetRootDir()

	var curr_brev_directories []string
	files.ReadJSON(rootDir+"/.brev/active_projects.json", &curr_brev_directories)

	for _, v := range curr_brev_directories {
		if (strings.Contains(cwd,v)) {
			return true
		}
	}
	return false
}

func CheckOutsideBrevErrorMessage() bool {

	if (IsInProjectDirectory()) {
		return true
	}

	rootDir := GetRootDir()
	var curr_brev_directories []string
	files.ReadJSON(rootDir+"/.brev/active_projects.json", &curr_brev_directories)

	// Exit with error message
	fmt.Println("Endpoint commands only work in a Brev project.")
	if (len(curr_brev_directories)==0) {
		// If no directories, check if they have some remote.

		// Get Projects
		token, _ := auth.GetToken()
		brevAgent := BrevAgent{
			Key: token,
		}
		raw_projects, _ := brevAgent.GetProjects()
		if (len(raw_projects)==0) {
			// Encourage them to create their first project
			fmt.Println("You haven't made a brev project yet! Try running 'brev init'")	

		} else {
			// Encourage them to pull one of their existing projects
			fmt.Println("Set up one of your existing projects.")
			fmt.Println("For example, run 'brev init " + raw_projects[0].Name + "'")
		}
		
	} else {
		// Print active brev projects
		fmt.Println("Active Brev projects on your computer: ")
		for _, v := range curr_brev_directories {
			fmt.Println("\t"+v)
		}
	}
	return false
}

func StringInList(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}