/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package package_project

import (
	"github.com/brevdev/brev-go-cli/internal/auth"
	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/brev_ctx"
	"github.com/brevdev/brev-go-cli/internal/terminal"
)

func addPackage(name string, t *terminal.Terminal) error {
	t.Vprint("\nAdding package " + t.Yellow(name))

	brevCtx, err := brev_ctx.New()
	if err != nil {
		return err
	}

	project, err := brevCtx.Local.GetProject()
	if err != nil {
		return err
	}

	_, err = brevCtx.Remote.SetPackage(*project, name)
	if err != nil {
		t.Errprintf(err, "\nFailed to add package %s", name)
		return err
	}

	t.Vprint(t.Green("\nPackage "))
	t.Vprint(t.Yellow("%s", name))
	t.Vprint(t.Green(" installed successfully 🥞"))

	return nil
}

func removePackage(name string, t *terminal.Terminal) error {
	t.Vprint("\nRemoving package " + t.Yellow(name))

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

	_, err = brevAgent.RemovePackage(packageToRemove.Id)
	if err != nil {
		t.Errprintf(err, "Failed to remove package %s", name)
		return err
	}

	t.Vprint(t.Green("\nPackage "))
	t.Vprint(t.Yellow("%s", name))
	t.Vprint(t.Green(" removed successfully 🥞"))

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
		if v.Status == "pending" {
			t.Vprintf("\t%s==%s %s\n", v.Name, v.Version, t.Yellow(v.Status))
		} else if v.Status == "installed" {
			t.Vprintf("\t%s==%s %s\n", v.Name, v.Version, t.Green(v.Status))
		} else {
			t.Vprintf("\t%s==%s %s\n", v.Name, v.Version, t.Red(v.Status))
		}
	}

	return nil
}
