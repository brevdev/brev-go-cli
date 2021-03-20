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
	"time"

	"github.com/brevdev/brev-go-cli/internal/auth"
	"github.com/brevdev/brev-go-cli/internal/brev_api"
	"github.com/brevdev/brev-go-cli/internal/brev_ctx"
	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
	"github.com/brevdev/brev-go-cli/internal/files"
	"github.com/brevdev/brev-go-cli/internal/requests"
	"github.com/briandowns/spinner"
)

func addEndpoint(name string, context *cmdcontext.Context) error {
	// Create endpoint
	
	localContext, err := brev_ctx.GetLocal()
	if (err != nil) {
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
	ep, err = brevAgent.CreateEndpoint(name, localContext.Project.Id)
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
		if (v.Name==name) {
			id=v.Id
		}
	}
	if (id=="") {
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
	fmt.Fprintln(context.VerboseOut, "Removing endpoint")
	// TODO: add the writer
	// spinner.WithWriter(os.Stderr)
	s := spinner.New(spinner.CharSets[39], 100*time.Millisecond)
	s.Start()
	_,err = brevAgent.RemoveEndpoint(id)
	if err != nil {
		context.PrintErr("", err)
		s.Stop()
		return err
	}
	s.Stop()
	
	fmt.Fprintf(context.VerboseOut, "Removed endpoint %s successfully",  name)

	// Remove the python file
	files.DeleteFile(name +".py")
	
	// Update the endpoints.json
	var updatedEndpoints []brev_api.Endpoint
	for _, v := range endpoints {
		if (v.Id!=id) {
			updatedEndpoints = append(updatedEndpoints, v)
		}
	}
	files.OverwriteJSON(endpointFilePath, updatedEndpoints)

	return nil
}

func runEndpoint(name string, method string, arg []string, jsonBody string, context *cmdcontext.Context) error {
	fmt.Fprintf(context.Out, "Run ep file %s %s %s", name, method, arg)

	localContext, err_dir := brev_ctx.GetLocal()
	if (err_dir != nil) {
		// handle this
		return err_dir
	}
	
	var endpoint brev_api.Endpoint
	for _, v := range localContext.Endpoints {
		if (v.Name == name) {
			endpoint = v
		}
	}

	fmt.Println(localContext.Project.Domain)
	fmt.Println(endpoint.Uri)


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
		Endpoint:    fmt.Sprintf("%s%s", localContext.Project.Domain, endpoint.Uri),
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
	err = rawResponse.DecodePayload(response)
	if err != nil {
		context.PrintErr("Failed to deserialize response payload", err)
		return err
	}

	

	fmt.Fprint(context.VerboseOut, "\n\n")
	fmt.Fprint(context.VerboseOut, rawResponse.StatusCode)

	jsonstr, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		context.PrintErr("Failed to serialize response", err)
		return err
	}

	fmt.Fprint(context.VerboseOut, jsonstr)
	fmt.Fprint(context.VerboseOut, string(jsonstr))

	// request2 := &requests.RESTRequest{
	// 	Method:      "POST",
	// 	Endpoint:    localContext.Project.Domain + endpoint.Uri,
	// 	QueryParams: params,
	// 	Headers: []requests.Header{
	// 		{Key: "Content-Type", Value: "application/json"},
	// 	},
	// }

	// rawResponse2, err := request2.Submit()
	// if err != nil {
	// 	context.PrintErr("Failed to run endpoint", err)
	// 	return err
	// }

	// var response2 map[string]string
	// err = rawResponse2.DecodePayload(&response2)
	// if err != nil {
	// 	context.PrintErr("Failed to deserialize response", err)
	// 	return err
	// }

	// fmt.Fprintln(context.Out, rawResponse2.StatusCode)
	// jsonstr2, err := json.MarshalIndent(response2, "", "  ")
	// if err != nil {
	// 	context.PrintErr("Failed to serialize response", err)
	// 	return err
	// }
	// fmt.Fprintln(context.Out, string(jsonstr2))

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
