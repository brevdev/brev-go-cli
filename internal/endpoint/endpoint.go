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
	"fmt"

	"github.com/brevdev/brev-go-cli/internal/requests"
)

func logic() {
	fmt.Println("endpoint called")
}

func add_endpoint(name string) {
	fmt.Printf("Create ep file %s", name)
}

func remove_endpoint(name string) {
	fmt.Printf("Remove ep file %s", name)
}

func run_endpoint(name string, method string, arg []string, jsonBody string) {
	fmt.Printf("Run ep file %s %s %s", name, method, arg)

    apiUrl := "https://dev-fjaq77pr.brev.dev/api/hi"

	response := requests.Get(requests.RequestParams{
		UrlString: apiUrl,
		QueryParams: arg,
	})

	fmt.Println(response.Status)
	fmt.Println(response.Body)

	response2 := requests.Post(requests.PostRequestParams{
		RequestParams: requests.RequestParams{
			UrlString: apiUrl,
			QueryParams: arg,
		},
		Body: "{'test':'value'}",
	})

	fmt.Println(response2.Status)
	fmt.Println(response2.Body)
	
}

func list_endpoints() {
	fmt.Println("List all endpoints")
}

func log_endpoint(name string) {
	fmt.Printf("Log ep file %s", name)
}
