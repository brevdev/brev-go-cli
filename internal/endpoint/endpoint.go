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
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
	"github.com/brevdev/brev-go-cli/internal/files"
	"github.com/brevdev/brev-go-cli/internal/requests"
)

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
	fmt.Fprintf(context.Out, "Remove ep file %s", name)

	return nil
}

func runEndpoint(name string, method string, arg []string, jsonBody string, context *cmdcontext.Context) error {
	fmt.Fprintf(context.Out, "Run ep file %s %s %s", name, method, arg)

	var params []requests.QueryParam
	for _, v := range arg {
		if strings.Contains(v, "=") {
			newArr := strings.Split(v, "=")
			params = append(params, requests.QueryParam{
				Key: newArr[0], Value: newArr[1]})
		}
	}

	request := &requests.RESTRequest{
		Method:      "GET",
		Endpoint:    "https://dev-fjaq77pr.brev.dev/api/hi",
		QueryParams: params,
		Headers: []requests.Header{
			{Key: "Content-Type", Value: "application/json"},
		},
	}
	rawResponse, err := request.Submit()
	if err != nil {
		context.PrintErr("Failed to run endpoint", err)
		return err
	}

	var response map[string]string
	err = rawResponse.DecodePayload(&response)
	if err != nil {
		context.PrintErr("Failed to deserialize response payload", err)
		return err
	}

	fmt.Fprint(context.Out, "\n\n")
	fmt.Fprint(context.Out, rawResponse.StatusCode)

	jsonstr, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		context.PrintErr("Failed to serialize response", err)
		return err
	}

	fmt.Fprint(context.Out, string(jsonstr))

	request2 := &requests.RESTRequest{
		Method:      "POST",
		Endpoint:    "https://dev-fjaq77pr.brev.dev/api/hi",
		QueryParams: params,
		Headers: []requests.Header{
			{Key: "Content-Type", Value: "application/json"},
		},
	}

	rawResponse2, err := request2.Submit()
	if err != nil {
		context.PrintErr("Failed to run endpoint", err)
		return err
	}

	var response2 map[string]string
	err = rawResponse2.DecodePayload(&response2)
	if err != nil {
		context.PrintErr("Failed to deserialize response", err)
		return err
	}

	fmt.Fprintln(context.Out, rawResponse2.StatusCode)
	jsonstr2, err := json.MarshalIndent(response2, "", "  ")
	if err != nil {
		context.PrintErr("Failed to serialize response", err)
		return err
	}
	fmt.Fprintln(context.Out, string(jsonstr2))

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
	for _, v := range endpointsResponse.Endpoints {
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
