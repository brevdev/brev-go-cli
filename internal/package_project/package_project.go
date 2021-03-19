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
)

func logic(context *cmdcontext.Context) error {
	fmt.Fprintln(context.Out, "package called")
	return nil
}


func addPackage(name string, context *cmdcontext.Context) error {

	localContext, err_dir := brev_ctx.GetLocal()
	if (err_dir != nil) {
		// handle this
		return err_dir
	}
	
	token, err := auth.GetToken()
	if err != nil {
		context.PrintErr("Failed to retrieve auth token", err)
		return err
	}
	brevAgent := brev_api.Agent{
		Key: token,
	}
	
	respPackage, err2 := brevAgent.AddPackage(localContext.Project.Id, name, context)
	if err2 != nil {
		context.PrintErr(fmt.Sprintf("Failed to add package %s", name), err2)
		return err2
	}

	fmt.Fprintln(context.Out, fmt.Sprintf("Package %s installed successfully.",&respPackage.Package.Name )) // this isn't working
	
	return nil
}