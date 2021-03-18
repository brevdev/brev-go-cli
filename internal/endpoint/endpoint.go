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
	"strings"

	"github.com/brevdev/brev-go-cli/internal/auth"
	"github.com/brevdev/brev-go-cli/internal/brev"
	"github.com/brevdev/brev-go-cli/internal/requests"
)

type Endpoint struct {
	Archived   bool     `json:"archived"`
	Code       string   `json:"code"`
	CreateDate string   `json:"create_date"`
	Id         string   `json:"id"`
	Methods    []string `json:"methods"`
	Name       string   `json:"name"`
	Namespace  string   `json:"namespace"`
	ProjectId  string   `json:"project_id"`
	Uri        string   `json:"uri"`
}

func add_endpoint(name string) {
	fmt.Printf("Create ep file %s", name)
	// TODO: create endpoint
	// TODO: get contents of .brev/endpoints.json
	// TODO: add new endpoint to .brev/endpoints.json
}

func remove_endpoint(name string) {
	fmt.Printf("Remove ep file %s", name)
}

func run_endpoint(name string, method string, arg []string, jsonBody string) {
	fmt.Printf("Run ep file %s %s %s", name, method, arg)

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
	raw_response, _ := request.Submit()
	var response map[string]string
	raw_response.DecodePayload(&response)
	fmt.Print("\n\n")
	fmt.Println(raw_response.StatusCode)
	jsonstr, _ := json.MarshalIndent(response, "", "  ")
	fmt.Println(string(jsonstr))

	request2 := &requests.RESTRequest{
		Method:      "POST",
		Endpoint:    "https://dev-fjaq77pr.brev.dev/api/hi",
		QueryParams: params,
		Headers: []requests.Header{
			{Key: "Content-Type", Value: "application/json"},
		},
	}
	raw_response2, _ := request2.Submit()
	var response2 map[string]string
	raw_response2.DecodePayload(&response2)
	fmt.Println(raw_response2.StatusCode)
	jsonstr2, _ := json.MarshalIndent(response2, "", "  ")
	fmt.Println(string(jsonstr2))

}

func list_endpoints() {
	token, _ := auth.GetToken()
	brevAgent := brev.BrevAgent{
		Key: token,
	}

	endpointsResponse, _ := brevAgent.GetEndpoints()
	// fmt.Println(endpointsResponse)
	var tempProj = endpointsResponse.Endpoints[0].ProjectId
	fmt.Printf("Endpoints in %s\n", tempProj)
	for _, v := range endpointsResponse.Endpoints {
		if (v.ProjectId==tempProj) {
			fmt.Printf("\tEp %s\n",v.Name)
			fmt.Printf("\t%s\n\n",v.Uri)

		}
	}

	// fmt.Println("List all endpoints")
	// var foo []Endpoint
	// files.ReadJSON("endpoints.json", &foo)
	// for _, v := range foo {
	// 	fmt.Println(v.Uri)
	// }
}

func log_endpoint(name string) {
	fmt.Printf("Log ep file %s", name)
}
