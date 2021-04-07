package endpoint

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/brev_ctx"
	"github.com/brevdev/brev-go-cli/internal/files"
	"github.com/brevdev/brev-go-cli/internal/requests"
	"github.com/brevdev/brev-go-cli/internal/terminal"
)

func addEndpoint(name string, t *terminal.Terminal) error {
	bar := t.NewProgressBar("\nAdding endpoint "+t.Yellow(name), func() {})

	brevCtx, err := brev_ctx.New()
	if err != nil {
		return err
	}

	// get current context project
	bar.Describe("Determining local project...")

	project, err := brevCtx.Local.GetProject()
	if err != nil {
		return err
	}

	// store endpoint in remote state
	bar.Describe("Submitting request to create new endpoint")
	endpoint, err := brevCtx.Remote.SetEndpoint(brev_api.Endpoint{
		ProjectId: project.Id,
		Name:      name,
	})
	if err != nil {
		return err
	}
	bar.AdvanceTo(30)

	// store endpoint in local state
	bar.Describe("Saving endpoint locally...")
	err = brevCtx.Local.SetEndpoint(*endpoint)
	if err != nil {
		return err
	}

	// create the endpoint code file
	cwd, err := os.Getwd()
	if err != nil {
		t.Errprint(err, "\nFailed to determine working directory")
		return err
	}

	err = files.OverwriteString(fmt.Sprintf("%s/%s.py", cwd, endpoint.Name), endpoint.Code)
	if err != nil {
		t.Errprint(err, "\nFailed to write endpoints to local file")
		return err
	}
	bar.AdvanceTo(100)

	t.Vprint(t.Green("\nEndpoint ") + t.Yellow("%s", name) + t.Green(" created and deployed 🥞"))

	return nil
}

func removeEndpoint(name string, t *terminal.Terminal) error {
	bar := t.NewProgressBar("Removing endpoint "+t.Yellow(name), func() {})
	bar.AdvanceTo(30)

	brevCtx, err := brev_ctx.New()
	if err != nil {
		return err
	}
	project, err := brevCtx.Local.GetProject()
	if err != nil {
		t.Errprint(err, "Cannot delete Endpoint. ")
		return err
	}
	eps, err := brevCtx.Remote.GetEndpoints(&brev_ctx.GetEndpointsOptions{
		Name: name,
	})
	if err != nil {
		t.Errprint(err, "Cannot delete Endpoint.")
		return err
	}
	if len(eps) == 0 {
		err := errors.New("endpoint doesn't exist")
		t.Errprint(err, "Cannot delete Endpoint.")
		return err
	}

	brevCtx.Remote.DeleteEndpoint(eps[0].Id)
	bar.Describe(t.Green("Endpoint ") + t.Yellow("%s", name) + t.Green(" deleted."))
	bar.AdvanceTo(60)
	// Remove the python file
	files.DeleteFile(name + ".py")

	// Update the endpoints.json
	allEndpoints, err := brevCtx.Remote.GetEndpoints(&brev_ctx.GetEndpointsOptions{
		ProjectID: project.Id,
	})
	if err != nil {
		t.Errprint(err, "Cannot delete Endpoint.")
		return err
	}
	files.OverwriteJSON(files.GetEndpointsPath(), allEndpoints)
	brevCtx.Local.SetEndpoints(allEndpoints)

	bar.Describe(t.Green("File ") + t.Yellow("%s.py", name) + t.Green(" removed."))
	bar.AdvanceTo(100)

	t.Vprint(t.Green("\nEndpoint ") + t.Yellow("%s", name) + t.Green(" removed from project ") + t.Yellow(project.Name) + " 🥞")

	return nil
}

func runEndpoint(name string, method string, arg []string, jsonBody string, t *terminal.Terminal) error {
	t.Vprint("\n")
	bar := t.NewProgressBar("Running endpoint "+t.Yellow(name), func() {})
	bar.AdvanceTo(40)

	brevCtx, err := brev_ctx.New()
	if err != nil {
		return err
	}

	// get local context project
	project, err := brevCtx.Local.GetProject()
	if err != nil {
		return err
	}

	bar.Describe("Preparing endpoint")
	bar.AdvanceTo(80)
	// get local endpoint for the given name
	endpoints, err := brevCtx.Local.GetEndpoints(&brev_ctx.GetEndpointsOptions{
		Name: name,
	})
	if err != nil {
		return err
	}
	if len(endpoints) != 1 {
		return fmt.Errorf(t.Red("unexpected number of endpoints: %d", len(endpoints)))
	}
	endpoint := endpoints[0]

	// prepare query params
	var params []requests.QueryParam
	for _, v := range arg {
		if strings.Contains(v, "=") {
			newArr := strings.Split(v, "=")
			params = append(params, requests.QueryParam{
				Key: newArr[0], Value: newArr[1]})
		}
	}

	// prepare payload
	var payload map[string]interface{}
	if jsonBody == "" {
		payload = nil
	} else if err := json.Unmarshal([]byte(jsonBody), &payload); err != nil {
		return fmt.Errorf(t.Red("failed to process JSON payload: %s", err))
	}

	bar.Describe("Submitting the request")
	bar.AdvanceTo(100)

	// submit request
	request := &requests.RESTRequest{
		Method:      method,
		Endpoint:    fmt.Sprintf("%s%s", project.Domain, endpoint.Uri),
		QueryParams: params,
		Headers: []requests.Header{
			{Key: "Content-Type", Value: "application/json"},
		},
		Payload: payload,
	}
	response, err := request.Submit()
	if err != nil {
		t.Errprint(err, "Failed to run endpoint")
		return err
	}

	// print output
	t.Vprint(t.Yellow("\n%s %s", request.Method, request.URI))
	if 200 <= response.StatusCode && response.StatusCode < 300 {
		t.Vprint(t.Green(" [%d]", response.StatusCode))
	} else if response.StatusCode >= 400 {
		t.Vprint(t.Red(" [%d]", response.StatusCode))
	} else {
		t.Vprint(t.Yellow(" [%d]", response.StatusCode))
	}

	jsonStr, err := response.PayloadAsPrettyJSONString()
	if err != nil {
		return err
	}

	t.Vprint("\n\nOutput:\n")
	t.Vprint(jsonStr)
	t.Vprint("\n\nLogs:\n")
	for _, header := range response.Headers {
		if header.Key == "x-stdout" {
			t.Vprint(header.Value)
		}
	}

	return nil
}

func listEndpoints(t *terminal.Terminal) error {
	brevCtx, err := brev_ctx.New()
	if err != nil {
		return err
	}

	// get current context project
	project, err := brevCtx.Local.GetProject()
	if err != nil {
		return err
	}

	// get remote project endpoints
	endpoints, err := brevCtx.Remote.GetEndpoints(&brev_ctx.GetEndpointsOptions{
		ProjectID: project.Id,
	})
	if err != nil {
		return err
	}

	// print
	t.Vprint(fmt.Sprintf("\nEndpoints in project %s:\n", project.Name))
	for _, endpoint := range endpoints {
		t.Vprint(fmt.Sprintf("\t%s:", t.Green(endpoint.Name)))
		t.Vprint(fmt.Sprintf("\t%s%s\n", project.Domain, endpoint.Uri))
	}

	return nil
}

func logEndpoint(name string, t *terminal.Terminal) error {
	t.Printf("Log ep file %s", name)
	return nil
}
