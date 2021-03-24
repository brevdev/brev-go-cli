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
package endpoint

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/brevdev/brev-go-cli/internal/auth"
	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/brev_ctx"
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
	"github.com/brevdev/brev-go-cli/internal/files"
	"github.com/brevdev/brev-go-cli/internal/requests"
)

func addEndpointV2(name string, context *cmdcontext.Context) error {
	brevCtx, err := brev_ctx.New()
	if err != nil {
		return err
	}

	// get current context project
	fmt.Fprint(context.Out, "Determining local project...\n")
	project, err := brevCtx.Local.GetProject()
	if err != nil {
		return err
	}
	fmt.Fprint(context.Out, fmt.Sprintf("Local project: %s\n", project.Name))

	// store endpoint in remote state
	fmt.Fprint(context.Out, "Submitting request to save new endpoint\n")
	endpoint, err := brevCtx.Remote.SetEndpoint(brev_api.Endpoint{
		ProjectId: project.Id,
		Name:      name,
	})
	if err != nil {
		return err
	}

	// store endpoint in local state
	fmt.Fprint(context.Out, "Saving endpoint locally...\n")
	err = brevCtx.Local.SetEndpoint(*endpoint)
	if err != nil {
		return err
	}

	// create the endpoint code file
	cwd, err := os.Getwd()
	if err != nil {
		context.PrintErr("Failed to determine working directory", err)
		return err
	}

	err = files.OverwriteString(fmt.Sprintf("%s/%s.py", cwd, endpoint.Name), endpoint.Code)
	if err != nil {
		context.PrintErr("Failed to write endpoints to local file", err)
		return err
	}

	return nil
}

func addEndpoint(name string, context *cmdcontext.Context) error {
	// Create endpoint
	proj, err := brev_api.GetActiveProject()
	if err != nil {
		context.PrintErr("Failed to get active project", err)
		return err
	}

	token, err := auth.GetToken()
	if err != nil {
		context.PrintErr("Failed to retrieve auth token", err)
		return err
	}
	brevAgent := brev_api.Agent{
		Key: token,
	}

	var ep *brev_api.ResponseUpdateEndpoint
	ep, err = brevAgent.CreateEndpoint(name, proj.Id)
	if err != nil {
		context.PrintErr("Failed to create endpoint", err)
		return err
	}

	fmt.Fprintln(context.VerboseOut, ep.Endpoint.Name+" created!")

	// Get contents of .brev/endpoints.json
	var allEps []brev_api.Endpoint
	err = files.ReadJSON(files.GetEndpointsPath(), &allEps)
	if err != nil {
		context.PrintErr("Failed to get endpoints", err)
		return err
	}

	// Add new endpoint to .brev/endpoints.json
	allEps = append(allEps, ep.Endpoint)
	err = files.OverwriteJSON(files.GetEndpointsPath(), allEps)
	if err != nil {
		context.PrintErr("Failed to write endpoints to .brev file", err)
		return err
	}

	// Create the endpoint code file
	cwd, err := os.Getwd()
	if err != nil {
		context.PrintErr("Failed to determine working directory", err)
		return err
	}

	err = files.OverwriteJSON(fmt.Sprintf("%s/%s.py", cwd, ep.Endpoint.Name), ep.Endpoint.Code)
	if err != nil {
		context.PrintErr("Failed to write endpoints to local file", err)
		return err
	}

	return nil
}

func removeEndpoint(name string, context *cmdcontext.Context) error {
	// Get the ID
	endpointFilePath := files.GetEndpointsPath()

	var endpoints []brev_api.Endpoint
	errFile := files.ReadJSON(endpointFilePath, &endpoints)
	if errFile != nil {
		return errFile
	}

	var id string
	for _, v := range endpoints {
		if v.Name == name {
			id = v.Id
		}
	}
	if id == "" {
		err := fmt.Errorf("Endpoint doesn't exist.")
		context.PrintErr("Cannot delete Endpoint. ", err)
		return err
	}

	// Remove the endpoint
	token, err := auth.GetToken()
	if err != nil {
		context.PrintErr("", err)
		return err
	}
	brevAgent := brev_api.Agent{
		Key: token,
	}

	// var ep *brev_api.ResponseRemoveEndpoint
	_, err = brevAgent.RemoveEndpoint(id)
	if err != nil {
		context.PrintErr("", err)
		return err
	}

	// Remove the python file
	files.DeleteFile(name + ".py")

	// Update the endpoints.json
	var updatedEndpoints []brev_api.Endpoint
	for _, v := range endpoints {
		if v.Id != id {
			updatedEndpoints = append(updatedEndpoints, v)
		}
	}
	files.OverwriteJSON(endpointFilePath, updatedEndpoints)

	return nil
}

func runEndpoint(name string, method string, arg []string, jsonBody string, context *cmdcontext.Context) error {
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
		return fmt.Errorf("unexpected number of endpoints: %d", len(endpoints))
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
		return fmt.Errorf("failed to process JSON payload: %s", err)
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
		context.PrintErr("Failed to run endpoint", err)
		return err
	}

	// print output
	fmt.Fprint(context.VerboseOut, fmt.Sprintf("%s %s [%d]", request.Method, request.URI, response.StatusCode))
	jsonStr, err := response.PayloadAsPrettyJSONString()
	if err != nil {
		return err
	}

	fmt.Fprint(context.VerboseOut, "\n\nOutput:\n")
	fmt.Fprint(context.VerboseOut, jsonStr)
	fmt.Fprint(context.VerboseOut, "\n\nLogs:\n")
	for _, header := range response.Headers {
		if header.Key == "x-stdout" {
			fmt.Fprint(context.VerboseOut, header)
		}
	}

	return nil
}

func listEndpointsV2(context *cmdcontext.Context) error {
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
	fmt.Fprintf(context.VerboseOut, "Endpoints in %s\n", project.Name)
	for _, endpoint := range endpoints {
		fmt.Fprintf(context.VerboseOut, "\tEp %s\n", endpoint.Name)
		fmt.Fprintf(context.VerboseOut, "\t%s\n\n", endpoint.Uri)
	}

	return nil
}

func listEndpoints(context *cmdcontext.Context) error {
	// get active project
	proj, err := brev_api.GetActiveProject()
	if err != nil {
		context.PrintErr("Failed to get active project", err)
		return err
	}

	token, err := auth.GetToken()
	if err != nil {
		context.PrintErr("Failed to retrieve auth token", err)
		return err
	}
	brevAgent := brev_api.Agent{
		Key: token,
	}

	endpointsResponse, err := brevAgent.GetEndpoints()
	if err != nil {
		context.PrintErr("Failed to get endpoints", err)
		return err
	}

	fmt.Fprintf(context.Out, "Endpoints in %s\n", proj.Name)
	for _, v := range endpointsResponse {
		if v.ProjectId == proj.Id {
			fmt.Fprintf(context.Out, "\tEp %s\n", v.Name)
			fmt.Fprintf(context.Out, "\t%s\n\n", v.Uri)
		}
	}
	return nil
}

func logEndpoint(name string, context *cmdcontext.Context) error {
	fmt.Fprintf(context.Out, "Log ep file %s", name)
	return nil
}
