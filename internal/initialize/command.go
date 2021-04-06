package initialize

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/brevdev/brev-go-cli/internal/auth"
	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/brev_ctx"
	"github.com/brevdev/brev-go-cli/internal/brev_errors"
	"github.com/brevdev/brev-go-cli/internal/files"
	"github.com/brevdev/brev-go-cli/internal/terminal"
)

func NewCmdInit(t *terminal.Terminal) *cobra.Command {
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

			bar0 := t.NewProgressBar("", 4, func() {})

			if project == "" {
				bar0.Describe(t.Yellow("\nInitializing new project"))
			} else {
				bar0.Describe(t.Yellow("\nInitializing project %s", project))
			}

			token, err := auth.GetToken()
			if err != nil {
				return err
			}
			bar0.Add(1)
			brevAgent := brev_api.Agent{
				Key: token,
			}
			projects, err := brevAgent.GetProjects()
			if err != nil {
				return fmt.Errorf("failed to retrieve projects %v", err)
			}
			bar0.Add(1)

			if project == "" {
				err = initNewProject(t)
				if err != nil {
					return err
				}
			}

			for _, v := range projects {

				if v.Name == project {
					err = initExistingProj(v, t)
					if err != nil {
						return fmt.Errorf("failed to initialize project %v", err)
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

func initExistingProj(project brev_api.Project, t *terminal.Terminal) error {

	cwd, err := os.Getwd()
	if err != nil {
		t.Errprint(err, "Failed to determine working directory")
		return err
	}

	bar1 := t.NewProgressBar("\nCloning Brev project in "+t.Yellow(cwd), 2, func() {})
	bar1.Describe("creating local files")

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
		t.Errprint(err, "Failed to write project to local file")
		return err
	}

	// Make endpoints.json
	err = files.OverwriteJSON(path+"/"+files.GetBrevDirectory()+"/"+files.GetEndpointsFile(), endpoints)
	if err != nil {
		t.Errprint(err, "Failed to write endpoints to local file")
		return err
	}

	// Create a global file with project directories
	var currBrevDirectories []string
	err = files.ReadJSON(files.GetActiveProjectsPath(), &currBrevDirectories)
	if err != nil {
		t.Errprint(err, "Failed to read projects directory")
		return err
	}

	if !brev_api.StringInList(path, currBrevDirectories) {
		currBrevDirectories = append(currBrevDirectories, path)
		err = files.OverwriteJSON(files.GetActiveProjectsPath(), currBrevDirectories)
		if err != nil {
			t.Errprint(err, "Failed to write projects to project file")
			return err
		}
	}

	bar1.Add(1)
	// TODO: copy shared code

	// Create endpoint files
	for _, v := range endpoints {
		err = files.OverwriteString(fmt.Sprintf("%s/%s.py", path, v.Name), v.Code)
		if err != nil {
			t.Errprint(err, "Failed to write code to local file")
			return err
		}
	}

	bar1.Describe(t.Green("\n\nBrev project %s cloned.", project.Name))
	bar1.Add(1)
	completionString := t.Yellow("\ncd %s", project.Name) + t.Green(" and get started!") + t.Green("\n\nHappy Hacking 🥞")

	t.Vprint(completionString)

	return nil
}

func initNewProject(t *terminal.Terminal) error {

	// Get Project Name (parent folder-- behavior just like git init)
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	bar1 := t.NewProgressBar("Creating Brev project in "+t.Yellow(cwd), 2, func() {})

	dirs := strings.Split(cwd, "/")
	projName := dirs[len(dirs)-1]

	// Create new project
	token, _ := auth.GetToken()
	brevAgent := brev_api.Agent{
		Key: token,
	}
	projectResponse, _ := brevAgent.CreateProject(projName)
	project := projectResponse.Project

	projectFilePath := cwd + "/" + files.GetBrevDirectory() + "/" + files.GetProjectsFile()
	endpointsFilePath := cwd + "/" + files.GetBrevDirectory() + "/" + files.GetEndpointsFile()
	activeProjectsFilePath := files.GetActiveProjectsPath()

	// Check if this is already an existing project
	if projectFileExists, err := files.Exists(projectFilePath); err != nil {
		return err
	} else if projectFileExists {
		return &brev_errors.InitExistingProjectFile{}
	}
	if endpointsFileExists, err := files.Exists(endpointsFilePath); err != nil {
		return err
	} else if endpointsFileExists {
		return &brev_errors.InitExistingEndpointsFile{}
	}

	bar1.Describe(t.Green("Creating local files..."))
	bar1.Add(1)

	// Make project.json
	err = files.OverwriteJSON(projectFilePath, project)
	if err != nil {
		t.Errprint(err, "Failed to write project to local file")
		return err
	}

	// Make endpoints.json
	err = files.OverwriteJSON(endpointsFilePath, []string{})
	if err != nil {
		t.Errprint(err, "Failed to write project to local file")
		return err
	}

	// Make active_projects.json if not exists
	if activeProjectsFilePathExists, err := files.Exists(activeProjectsFilePath); err != nil {
		return err
	} else if !activeProjectsFilePathExists {
		err = files.OverwriteJSON(activeProjectsFilePath, []string{})
		if err != nil {
			return fmt.Errorf("failed to write active projects to global file %v", err)
		}
	}

	// TODO: create shared code module

	// Add to path
	var currBrevDirectories []string
	err = files.ReadJSON(files.GetActiveProjectsPath(), &currBrevDirectories)
	if err != nil {
		t.Errprint(err, "Failed to read projects directory")
		return err
	}

	if !brev_api.StringInList(cwd, currBrevDirectories) {
		currBrevDirectories = append(currBrevDirectories, cwd)
		err = files.OverwriteJSON(activeProjectsFilePath, currBrevDirectories)
		if err != nil {
			t.Errprint(err, "Failed to write projects to project file")
			return err
		}
	}

	bar1.Describe(t.Green("\n\nBrev project %s created and deployed.", projName))
	bar1.Add(1)
	completionString := t.Green(t.Yellow("\ncd %s", projName) + t.Green(" and get started!") + t.Green("\n\nHappy Hacking 🥞"))
	t.Vprint(completionString)

	return nil

}
