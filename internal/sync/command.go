
package sync

import (
	"fmt"
	"os"
	"strings"

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
