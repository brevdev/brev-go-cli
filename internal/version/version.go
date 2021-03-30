package version

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/brevdev/brev-go-cli/internal/config"
	"github.com/brevdev/brev-go-cli/internal/requests"
	"github.com/brevdev/brev-go-cli/internal/terminal"
)

const (
	cliReleaseURL = "https://api.github.com/repos/brevdev/brev-go-cli/releases/latest"
)

var green = color.New(color.FgGreen).SprintfFunc()

var upToDateString = `
Current version: %s

` + green("You're up to date!")

var outOfDateString = `
Current version: %s

` + green("A new version of brev has been released!") + `

Version: %s

Details: %s

` + green("run 'brew upgrade brevdev/tap/brev' to upgrade") + `

%s
`

type githubReleaseMetadata struct {
	TagName      string `json:"tag_name"`
	IsDraft      bool   `json:"draft"`
	IsPrerelease bool   `json:"prerelease"`
	Name         string `json:"name"`
	Body         string `json:"body"`
}

func BuildVersionString(t *terminal.Terminal) (string, error) {
	currentVersion := config.GetVersion()

	githubRelease, err := getLatestGithubReleaseMetadata()
	if err != nil {
		t.Errprint(err, "Failed to retrieve latest version")
		return "", err
	}

	var versionString string
	if githubRelease.TagName == currentVersion {
		versionString = fmt.Sprintf(
			upToDateString,
			currentVersion,
		)
	} else {
		versionString = fmt.Sprintf(
			outOfDateString,
			currentVersion,
			githubRelease.TagName,
			githubRelease.Name,
			githubRelease.Body,
		)
	}
	return versionString, nil
}

func getLatestGithubReleaseMetadata() (*githubReleaseMetadata, error) {
	request := &requests.RESTRequest{
		Method:   "GET",
		Endpoint: cliReleaseURL,
	}
	response, err := request.Submit()
	if err != nil {
		return nil, err
	}

	var payload githubReleaseMetadata
	err = response.UnmarshalPayload(&payload)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}
