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
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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

	for i := 0; i < len(arg); i++ {
		fmt.Printf(arg[i])
	}


    apiUrl := "https://dev-fjaq77pr.brev.dev/api/hi"
    data := url.Values{}
    data.Set("name", "foo")
    data.Set("surname", "bar")

    // u, _ := url.ParseRequestURI(apiUrl)
    // urlStr := u.String() // "https://api.com/user/"

    client := &http.Client{}

	var body io.Reader;
	
	if method=="GET" {
		fmt.Println("YYYEEEPP")
		body = nil
	} else {
		body = strings.NewReader(data.Encode())
	}

    r, _ := http.NewRequest("GET", apiUrl, body) // URL-encoded payload
    // r, _ := http.NewRequest("GET", apiUrl, body) // URL-encoded payload
    r.Header.Add("Content-Type", "application/json")

    resp, _ := client.Do(r)
    fmt.Println(resp.Status)
	
	resp_body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(resp_body))

}

func list_endpoints() {
	fmt.Println("List all endpoints")
}

func log_endpoint(name string) {
	fmt.Printf("Log ep file %s", name)
}
