package endpoint

import (
	"encoding/json"
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
	t.Vprint("\nAdding endpoint " + t.Yellow(name))

	brevCtx, err := brev_ctx.New()
	if err != nil {
		return err
	}

	// get current context project
	t.Print("Determining local project...\n")
	project, err := brevCtx.Local.GetProject()
	if err != nil {
		return err
	}
	t.Print(fmt.Sprintf("Local project: %s\n", project.Name))

	// store endpoint in remote state
	t.Print("Submitting request to save new endpoint\n")
	endpoint, err := brevCtx.Remote.SetEndpoint(brev_api.Endpoint{
		ProjectId: project.Id,
		Name:      name,
	})
	if err != nil {
		return err
	}
	t.Vprint(t.Green("\nEndpoint "))
	t.Vprint(t.Yellow("%s", name))
	t.Vprint(t.Green(" created and deployed 🚀"))

	// store endpoint in local state
	t.Print("Saving endpoint locally...\n")
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

	err = files.OverwriteString(fmt.Sprintf("\n%s/%s.py", cwd, endpoint.Name), endpoint.Code)
	if err != nil {
		t.Errprint(err, "\nFailed to write endpoints to local file")
		return err
	}
	t.Vprint(t.Yellow("\n%s.py", name))
	t.Vprint(t.Green(" created 🥞"))

	return nil
}

func removeEndpoint(name string, t *terminal.Terminal) error {
	t.Vprint("\nRemoving endpoint " + t.Yellow(name))

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

	brevCtx.Remote.DeleteEndpoint(eps[0].Id)
	t.Vprint(t.Green("\nEndpoint "))
	t.Vprint(t.Yellow("%s", name))
	t.Vprint(t.Green(" deleted."))

	// Remove the python file
	files.DeleteFile(name + ".py")

	// Update the endpoints.json
	allEndpoints, err := brevCtx.Remote.GetEndpoints(&brev_ctx.GetEndpointsOptions{
		ProjectID: project.Id,
	})
	for _, v := range allEndpoints {
		fmt.Println(v.Name)
	}
	if err != nil {
		t.Errprint(err, "Cannot delete Endpoint.")
		return err
	}
	files.OverwriteJSON(files.GetEndpointsPath(), allEndpoints)
	brevCtx.Local.SetEndpoints(allEndpoints)

	t.Vprint(t.Green("\nFile "))
	t.Vprint(t.Yellow("%s.py", name))
	t.Vprint(t.Green(" removed."))

	return nil
}

func runEndpoint(name string, method string, arg []string, jsonBody string, t *terminal.Terminal) error {
	t.Vprint("\nRunning endpoint " + t.Yellow(name))

	brevCtx, err := brev_ctx.New()
	if err != nil {
		return err
	}

	// get local context project
	project, err := brevCtx.Local.GetProject()
	if err != nil {
		return err
	}

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

	// print
	t.Vprintf("Endpoints in project %s:\n", project.Name)
	for _, endpoint := range endpoints {
		t.Vprintf("\t%s:\n", t.Green(endpoint.Name))
		t.Vprintf("\t%s%s\n\n", project.Domain, endpoint.Uri)
	}

	return nil
}

func logEndpoint(name string, t *terminal.Terminal) error {
	t.Printf("Log ep file %s", name)
	return nil
}
