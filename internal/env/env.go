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
package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/brevdev/brev-go-cli/internal/auth"
	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/brev_ctx"
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
)

func logic(context *cmdcontext.Context) error {
	return nil
}

func addVariable(name string, context *cmdcontext.Context) error {
	
	localContext, err_dir := brev_ctx.GetLocal()
	if (err_dir != nil) {
		// handle this
		return err_dir
	}
	
	fmt.Fprintf(context.VerboseOut, "Enter value for %s: ", name)	
	reader := bufio.NewReader(os.Stdin)
	
	value, _ := reader.ReadString('\n')
    // convert CRLF to LF
    value = strings.Replace(value, "\n", "", -1)

	token, err := auth.GetToken()
	if err != nil {
		context.PrintErr("Failed to retrieve auth token", err)
		return err
	}
	brevAgent := brev_api.Agent{
		Key: token,
	}

	brevAgent.AddVariable(localContext.Project.Id, name, value)
	
	fmt.Fprintf(context.VerboseOut, "\nVariable %s added to your project.", name)	

	return nil
}


func removeVariable(name string, context *cmdcontext.Context) error {
	
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

	// Find ID of variable requested to delete
	projVars, err := brevAgent.GetVariables(localContext.Project.Id)
	if (err != nil) {
		context.PrintErr("Couldn't access project environment variables.", err)
		return err
	}

	var varId string
	for _, v := range projVars {
		if v.Name==name {
			varId = v.Id
		}
	}
	if varId=="" {
		context.PrintErr(fmt.Sprintf("There isn't a variable in your project named %s.", name), err)
		return err
	}

	// Remove variable by ID
	_,err = brevAgent.RemoveVariable(varId)
	if err != nil {
		context.PrintErr("Couldn't remove the variable.", err)
		return err
	}
	
	fmt.Fprintf(context.VerboseOut, "\nVariable %s removed from your project.", name)	

	return nil
}
