package package_project

import (
	"fmt"

	"github.com/brevdev/brev-go-cli/internal/auth"
	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/brev_ctx"
	"github.com/brevdev/brev-go-cli/internal/terminal"
)

func addPackage(name string, t *terminal.Terminal) error {
	bar := t.NewProgressBar("Adding package "+t.Yellow(name), func() {})

	brevCtx, err := brev_ctx.New()
	if err != nil {
		return err
	}

	bar.AdvanceTo(40)
	project, err := brevCtx.Local.GetProject()
	if err != nil {
		return err
	}

	_, err = brevCtx.Remote.SetPackage(*project, name)
	if err != nil {
		t.Errprintf(err, "\nFailed to add package %s", name)
		return err
	}

	bar.AdvanceTo(100)
	finalStr := t.Green("Package ") + t.Yellow("%s", name) + t.Green(" added successfully 🥞")
	t.Vprintln(finalStr)

	t.Vprintln(t.Yellow(`

The package might take a moment to fully install.
'brev package list' to see the status of all packages.
	`))

	return nil
}

func removePackage(name string, t *terminal.Terminal) error {
	bar := t.NewProgressBar("Removing package "+t.Yellow(name), func() {})

	packages, err := GetPackages(t)
	if err != nil {
		return nil
	}

	var packageToRemove brev_api.ProjectPackage
	for _, v := range packages {
		if v.Name == name {
			packageToRemove = v
		}
	}

	token, err := auth.GetToken()
	if err != nil {
		t.Errprint(err, "Failed to retrieve auth token")
		return err
	}
	brevAgent := brev_api.Agent{
		Key: token,
	}

	bar.AdvanceTo(60)
	_, err = brevAgent.RemovePackage(packageToRemove.Id)
	if err != nil {
		t.Errprintf(err, "Failed to remove package %s", name)
		return err
	}
	finalStr := t.Green("\nPackage ") + t.Yellow("%s", name) + t.Green(" removed successfully 🥞")

	// bar.Describe(finalStr)
	bar.AdvanceTo(100)

	t.Vprintln(finalStr)

	return nil
}

func listPackages(t *terminal.Terminal) error {
	packages, err := GetPackages(t)
	if err != nil {
		return nil
	}

	brevCtx, err := brev_ctx.New()
	if err != nil {
		return err
	}

	project, err := brevCtx.Local.GetProject()
	if err != nil {
		return err
	}

	t.Vprintf("Packages installed on project %s:\n", project.Name)

	for _, v := range packages {
		installStr := fmt.Sprintf("\t%s==%s ", v.Name, v.Version)
		if v.Status == "pending" {
			t.Vprintln(installStr + t.Yellow("%s", v.Status))
		} else if v.Status == "installed" {
			t.Vprintln(installStr + t.Green("%s", v.Status))
		} else {
			t.Vprintln(installStr + t.Red("%s", v.Status))
		}
	}

	return nil
}
