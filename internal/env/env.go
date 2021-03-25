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
	"errors"
	"fmt"
	"os"
	"syscall"

	"golang.org/x/term"

	"github.com/brevdev/brev-go-cli/internal/auth"
	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/brev_ctx"
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
)

func logic(context *cmdcontext.Context) error {
	return nil
}

func addVariable(name string, context *cmdcontext.Context) error {

	brevCtx, err := brev_ctx.New()
	if err != nil {
		return err
	}

	project, err := brevCtx.Local.GetProject()
	if err != nil {
		return err
	}

	fmt.Fprintf(context.VerboseOut, "Enter value for %s: ", name)

	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		os.Exit(1)
	}
	value := string(bytepw)

	brevCtx.Remote.SetVariable(*project, name, value)

	fmt.Fprintf(context.VerboseOut, "\nVariable %s added to your project.", name)

	return nil
}

func removeVariable(name string, context *cmdcontext.Context) error {

	brevCtx, err := brev_ctx.New()
	if err != nil {
		return err
	}

	project, err := brevCtx.Local.GetProject()
	if err != nil {
		return err
	}

	projVars, err := brevCtx.Remote.GetVariables(*project, &brev_ctx.GetVariablesOptions{
		Name: name,
	})
	if err != nil {
		return errors.New(fmt.Sprintf("There isn't a variable in your project named %s.", name))
	}

	token, err := auth.GetToken()
	if err != nil {
		context.PrintErr("Failed to retrieve auth token", err)
		return err
	}
	brevAgent := brev_api.Agent{
		Key: token,
	}

	// Remove variable by ID
	_, err = brevAgent.RemoveVariable(projVars[0].Id)
	if err != nil {
		context.PrintErr("Couldn't remove the variable.", err)
		return err
	}

	fmt.Fprintf(context.VerboseOut, "\nVariable %s removed from your project.", name)

	return nil
}
