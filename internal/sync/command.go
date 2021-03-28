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
package sync

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/brevdev/brev-go-cli/internal/auth"
	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/brev_ctx"
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
	"github.com/brevdev/brev-go-cli/internal/files"
)

func push(context *cmdcontext.Context) error {

	// TODO: push module/shared code

	token, err := auth.GetToken()
	if err != nil {
		context.PrintErr("", err)
		return err
	}
	brevAgent := brev_api.Agent{
		Key: token,
	}

	brevCtx, err := brev_ctx.New()
	if err != nil {
		return err
	}

	project, err := brevCtx.Local.GetProject()
	if err != nil {
		return err
	}

	endpoints, err := brevCtx.Local.GetEndpoints(&brev_ctx.GetEndpointsOptions{
		ProjectID: project.Id,
	})

	if err != nil {
		return err
	}

	for _, v := range endpoints {
		fmt.Fprintf(context.VerboseOut, "\nUpdating ep %s", v.Name)

		brevAgent.UpdateEndpoint(v.Id, brev_api.RequestUpdateEndpoint{
			Name:    v.Name,
			Methods: v.Methods,
			Code:    v.Code,
		})
	}

	return nil
}

func pull(context *cmdcontext.Context) error {

	// TODO: module/shared code

	brevCtx, err := brev_ctx.New()
	if err != nil {
		return err
	}

	project, err := brevCtx.Local.GetProject()
	if err != nil {
		return err
	}

	remoteEndpoints, err := brevCtx.Remote.GetEndpoints(&brev_ctx.GetEndpointsOptions{
		ProjectID: project.Id,
	})
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		context.PrintErr("Failed to determine working directory", err)
		return err
	}

	paths, err := brevCtx.Global.GetProjectPaths()
	if err != nil {
		return err
	}

	var path string
	for _, v := range paths {
		if strings.Contains(cwd, v) {
			path = v
		}
	}
	if path == "" {
		return errors.New("this is not a Brev directory")
	}

	for _, v := range remoteEndpoints {
		fmt.Fprintf(context.VerboseOut, "\nPulling ep %s", v.Name)

		err = files.OverwriteString(fmt.Sprintf("%s/%s.py", path, v.Name), v.Code)
		if err != nil {
			context.PrintErr("Failed to write code to local file", err)
			return err
		}
	}

	brevCtx.Local.SetEndpoints(remoteEndpoints)

	return nil
}
