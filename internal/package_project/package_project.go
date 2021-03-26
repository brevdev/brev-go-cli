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
	"fmt"

	"github.com/brevdev/brev-go-cli/internal/auth"
	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/brev_ctx"
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
	"github.com/fatih/color"
)

func logic(context *cmdcontext.Context) error {
	fmt.Fprintln(context.Out, "package called")
	return nil
}

func addPackage(name string, context *cmdcontext.Context) error {
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
		context.PrintErr(fmt.Sprintf("Failed to add package %s", name), err)
		return err
	}

	fmt.Fprintf(context.VerboseOut, "Package %s installed successfully.", name) // this isn't working

	return nil
}

func removePackage(name string, context *cmdcontext.Context) error {

	packages, err := GetPackages(context)
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
		context.PrintErr("Failed to retrieve auth token", err)
		return err
	}
	brevAgent := brev_api.Agent{
		Key: token,
	}

	_, err = brevAgent.RemovePackage(packageToRemove.Id)
	if err != nil {
		context.PrintErr(fmt.Sprintf("Failed to remove package %s", name), err)
		return err
	}

	fmt.Fprintf(context.VerboseOut, "\nPackage %s removed successfully.", packageToRemove.Name) // this isn't working

	return nil
}

func listPackages(context *cmdcontext.Context) error {

	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	packages, err := GetPackages(context)
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

	fmt.Fprintf(context.VerboseOut, "Packages installed on project %s:\n", project.Name)

	for _, v := range packages {
		if v.Status == "pending" {
			fmt.Fprintf(context.VerboseOut, "\t%s==%s %s\n", v.Name, v.Version, yellow(v.Status))
		} else if v.Status == "installed" {
			fmt.Fprintf(context.VerboseOut, "\t%s==%s %s\n", v.Name, v.Version, green(v.Status))
		} else {
			fmt.Fprintf(context.VerboseOut, "\t%s==%s %s\n", v.Name, v.Version, red(v.Status))

		}
	}

	return nil
}
