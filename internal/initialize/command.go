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

func NewCmdClone(t *terminal.Terminal) *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:         "clone",
		Short:       "Clone a Brev Project",
		Annotations: map[string]string{"project": ""},
		Long:        "Clone an existing Brev project",
		Example: `  // To clone your existing Brev project
  brev clone --name your_project_name`,
		RunE: func(cmd *cobra.Command, args []string) error {

			bar := t.NewProgressBar("", func() {})

			bar.Describe(t.Yellow("Cloning project %s", name))

			token, err := auth.GetToken()
			if err != nil {
				return err
			}
			bar.AdvanceTo(20)
			brevAgent := brev_api.Agent{
				Key: token,
			}
			projects, err := brevAgent.GetProjects()
			if err != nil {
				return fmt.Errorf("failed to retrieve projects %v", err)
			}
			bar.AdvanceTo(30)

			if name == "" {
				err = initNewProject(t, bar)
				if err != nil {
					return err
				}
			}

			for _, v := range projects {

				if v.Name == name {
					err = initExistingProj(v, t, bar)
					if err != nil {
						return fmt.Errorf("failed to initialize project %v", err)
					}
					break // in case of error where multiple projects share name. We should prohibit this.
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&name, "name", "p", "", "Project Name")
	cmd.MarkFlagRequired("name")
	cmd.RegisterFlagCompletionFunc("name", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getProjectNames(), cobra.ShellCompDirectiveNoSpace
	})

	return cmd
}

func NewCmdInit(t *terminal.Terminal) *cobra.Command {

	cmd := &cobra.Command{
		Use:         "init",
		Annotations: map[string]string{"project": ""},
		Short:       "Initialize a Brev Project",
		Long:        "Initialize a Brev project.",
		Example: `  // To init new project in current directory
  brev init`,
		RunE: func(cmd *cobra.Command, args []string) error {

			bar := t.NewProgressBar("", func() {})

			bar.Describe(t.Yellow("Initializing new project"))

			err := initNewProject(t, bar)
			if err != nil {
				return err
			}

			return nil
		},
	}

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

func initExistingProj(project brev_api.Project, t *terminal.Terminal, bar *terminal.ProgressBar) error {

	cwd, err := os.Getwd()
	if err != nil {
		t.Errprint(err, "Failed to determine working directory")
		return err
	}

	bar.Describe("\nCloning Brev project in " + t.Yellow(cwd))
	bar.Describe("creating local files")

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

	// Create endpoint files
	for _, v := range endpoints {
		err = files.OverwriteString(fmt.Sprintf("%s/%s.py", path, v.Name), v.Code)
		if err != nil {
			t.Errprint(err, "Failed to write code to local file")
			return err
		}
	}

	// Create shared code/module
	module, err := brevCtx.Remote.GetModule(&brev_ctx.GetModulesOptions{
		ProjectID: project.Id,
	})
	if err != nil {
		return err
	}

	err = files.OverwriteString(fmt.Sprintf("%s/%s.py", path, module.Name), module.Source)
	if err != nil {
		return err
	}
	bar.Describe(t.Green("Brev project %s cloned.", project.Name))
	bar.AdvanceTo(100)
	completionString := t.Yellow("\ncd %s", project.Name) + t.Green(" and get started!") + t.Green("\n\nHappy Hacking 🥞")

	t.Vprint(completionString)

	return nil
}

func initNewProject(t *terminal.Terminal, bar *terminal.ProgressBar) error {

	// Get Project Name (parent folder-- behavior just like git init)
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	bar.Describe("Creating Brev project in " + t.Yellow(cwd))

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

	bar.Describe(t.Green("Creating local files..."))
	bar.AdvanceTo(40)

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

	// Create shared code/module
	brevCtx, err := brev_ctx.New()
	if err != nil {
		return err
	}
	proj, err := brevCtx.Local.GetProject()
	if err != nil {
		return err
	}
	module, err := brevCtx.Remote.GetModule(&brev_ctx.GetModulesOptions{
		ProjectID: proj.Id,
	})
	if err != nil {
		return err
	}
	err = files.OverwriteString(fmt.Sprintf("%s/%s.py", cwd, module.Name), module.Source)
	if err != nil {
		return err
	}

	bar.Describe(t.Green("Brev project %s created and deployed.", projName))
	bar.AdvanceTo(100)
	completionString := t.Green("\n\nHappy Hacking 🥞")
	t.Vprint(completionString)

	return nil

}
