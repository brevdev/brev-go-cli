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

	path, err := getRootProjectDir(t)
	if err != nil {
		return err
	}

	brevCtx, err := brev_ctx.New()
	if err != nil {
		return err
	}

	project, err := brevCtx.Local.GetProject()
	if err != nil {
		return err
	}

	// update module
	module, err := brevCtx.Remote.GetModule(&brev_ctx.GetModulesOptions{ProjectID: project.Id})
	if err != nil {
		return err
	}
	module.Source, err = files.ReadString(fmt.Sprintf("%s/%s.py", path, module.Name))
	if err != nil {
		return err
	}

	_, err = brevCtx.Remote.SetModule(&brev_ctx.SetModulesOptions{
		ProjectID: project.Id,
		ModuleID:  module.Id,
		Source:    module.Source,
	})
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

	module, err := brevCtx.Remote.GetModule(&brev_ctx.GetModulesOptions{
		ProjectID: project.Id,
	})
	if err != nil {
		return err
	}
	t.Vprint(t.Green("\nPulling %s", module.Name))

	err = files.OverwriteString(fmt.Sprintf("%s/%s.py", path, module.Name), module.Source)

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

func diffCmd(t *terminal.Terminal) error {

	numChanges := 0

	brevCtx, err := brev_ctx.New()
	if err != nil {
		return err
	}

	project, err := brevCtx.Local.GetProject()
	if err != nil {
		return err
	}
	localEps, err := brevCtx.Local.GetEndpoints(&brev_ctx.GetEndpointsOptions{
		ProjectID: project.Id,
	})
	var localEPIds []string
	for _, v := range localEps {
		localEPIds = append(localEPIds, v.Id)
	}
	if err != nil {
		return err
	}
	remoteEps, err := brevCtx.Remote.GetEndpoints(&brev_ctx.GetEndpointsOptions{
		ProjectID: project.Id,
	})
	var remoteEPIds []string
	remoteEpMap := make(map[string]brev_api.Endpoint)
	for _, v := range remoteEps {
		remoteEPIds = append(remoteEPIds, v.Id)
		remoteEpMap[v.Id] = v
	}
	if err != nil {
		return err
	}
	t.Vprint(t.Yellow("Diff for Project %s :", project.Name))

	// per local endpoint, diff the remote contents
	for _, v := range localEps {
		// if the local ep has a remote counter part, run a diff
		if brev_api.StringInList(v.Id, remoteEPIds) {
			path, err := getRootProjectDir(t)
			if err != nil {
				return err
			}

			v.Code, err = files.ReadString(fmt.Sprintf("%s/%s.py", path, v.Name))
			if err != nil {
				return err
			}
			diff := diffTwoFiles(remoteEpMap[v.Id].Code, v.Code)
			diffString := printDiff(v.Name, diff, t)
			if len(diffString) > 0 {
				t.Vprint(diffString)
				numChanges += 1
			}
		} else {
			// The endpoint doesn't exist in remote
			diff := diffTwoFiles("", v.Code)
			diffString := printDiff(v.Name, diff, t)
			if len(diffString) > 0 {
				t.Vprint(diffString)
				numChanges += 1
			}
		}
	}
	// if remote endpoint isn't local, then it needs to be pulled
	for _, v := range remoteEps {
		if !brev_api.StringInList(v.Id, localEPIds) {
			diff := diffTwoFiles(remoteEpMap[v.Id].Code, "")
			diffString := printDiff(v.Name, diff, t)
			if len(diffString) > 0 {
				t.Vprint(diffString)
				numChanges += 1
			}
		}
	}

	if numChanges == 0 {
		t.Vprint(t.Green("All Synced 🥞"))
	}

	return nil
}

func diffTwoFiles(s1 string, s2 string) string {
	s1Trimmed := strings.TrimSpace(s1)
	s2Trimmed := strings.TrimSpace(s2)
	return diff.LineDiff(s1Trimmed, s2Trimmed)

}

func printDiff(filename string, diff string, t *terminal.Terminal) string {

	diffOutputString := ""
	totalDiffLines := 0
	for _, v := range strings.Split(diff, "\n") {
		if strings.Compare(string(v[0]), "+") == 0 {
			diffOutputString += "\n" + t.Green(v)
			totalDiffLines += 1
		} else if strings.Compare(string(v[0]), "-") == 0 {
			diffOutputString += "\n" + t.Red(v)
			totalDiffLines += 1
		}
	}
	if totalDiffLines > 0 {
		diffOutputString = t.Yellow("%s.py: ", filename) + diffOutputString + "\n"
	}
	return diffOutputString
}
