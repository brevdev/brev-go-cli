package version

import (
	"fmt"



	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
	"github.com/brevdev/brev-go-cli/internal/config"
	"github.com/brevdev/brev-go-cli/internal/requests"
)

const (
	cliReleaseURL = "https://api.github.com/repos/brevdev/brev-go-cli/releases/latest"
)

var upToDateString = `
Current version: %s

You're up to date!
`

var outOfDateString = `
Current version: %s

A new version of brev has been released!

Version: %s

Details: %s

%s
`

type githubReleaseMetadata struct {
	TagName      string `json:"tag_name"`
	IsDraft      bool   `json:"draft"`
	IsPrerelease bool   `json:"prerelease"`
	Name         string `json:"name"`
	Body         string `json:"body"`
}

func buildVersionString(context *cmdcontext.Context) (string, error) {
	currentVersion := config.GetVersion()

	githubRelease, err := getLatestGithubReleaseMetadata()
	if err != nil {
		context.PrintErr("Failed to retrieve latest version", err)
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
	err = response.DecodePayload(&payload)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}
