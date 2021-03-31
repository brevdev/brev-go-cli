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
	"fmt"
	"os"
	"strings"

	"github.com/andreyvit/diff"
	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/brev_ctx"
	"github.com/brevdev/brev-go-cli/internal/files"
	"github.com/brevdev/brev-go-cli/internal/terminal"
)

func push(t *terminal.Terminal) error {

	// TODO: push module/shared code
	t.Vprint(t.Green("\nPushing your changes..."))

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
		t.Vprint(t.Green("\nUpdating ep %s", v.Name))

		path, err := getRootProjectDir(t)
		if err != nil {
			return err
		}

		v.Code, err = files.ReadString(fmt.Sprintf("%s/%s.py", path, v.Name))
		if err != nil {
			return err
		}

		brevCtx.Remote.SetEndpoint(brev_api.Endpoint{
			Id:      v.Id,
			Name:    v.Name,
			Methods: v.Methods,
			Code:    v.Code,
		})

	}

	t.Vprint(t.Green("\n\nYour project is synced 🥞"))

	return nil
}

func pull(t *terminal.Terminal) error {

	// TODO: module/shared code
	t.Vprint(t.Green("\nPulling changes from the console..."))

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

	path, err := getRootProjectDir(t)
	if err != nil {
		return err
	}

	for _, v := range remoteEndpoints {
		t.Vprint(t.Green("\nPulling ep %s", v.Name))

		err = files.OverwriteString(fmt.Sprintf("%s/%s.py", path, v.Name), v.Code)
		if err != nil {
			t.Errprint(err, "Failed to write code to local file")
			return err
		}
	}

	brevCtx.Local.SetEndpoints(remoteEndpoints)

	t.Vprint(t.Green("\n\nYour project is synced 🥞"))

	return nil
}

func getRootProjectDir(t *terminal.Terminal) (string, error) {

	brevCtx, err := brev_ctx.New()
	if err != nil {
		return "", err
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Errprint(err, "Failed to determine working directory")
		return "", err
	}

	paths, err := brevCtx.Global.GetProjectPaths()
	if err != nil {
		return "", err
	}

	var path string
	for _, v := range paths {
		if strings.Contains(cwd, v) {
			path = v
		}
	}
	return path, nil
}

func diffCmd() error {
	green := color.New(color.FgGreen).SprintfFunc()
	red := color.New(color.FgRed).SprintfFunc()
	// Per endpoint,
	// perform diff
	var s1 = `
	test test this is a test
	new line in test

	hello
	`
	var s2 = `
	test test this is a test

	hello
	`

	diff := diffTwoFiles(s1, s2)

	// fmt.Println(diff)

	for _, v := range strings.Split(diff, "\n") {
		if strings.Compare(string(v[0]), "+") == 0 {
			fmt.Println(green(v))
		} else if strings.Compare(string(v[0]), "-") == 0 {
			fmt.Println(red(v))
		}
	}

	return nil
}

func diffTwoFiles(s1 string, s2 string) string {
	s1Trimmed := strings.TrimSpace(s1)
	s2Trimmed := strings.TrimSpace(s2)
	return diff.LineDiff(s1Trimmed, s2Trimmed)

}
