package requests

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// EXAMPLE USAGE:
/*
	// GET
	response := requests.Get(requests.RequestParams{
		UrlString: "https://wqgtqgir.brev.dev/api/CLI2",
		QueryParams: []string {"hi=bye", "hiii=byyyeee"},
	})

	fmt.Println(response.Status)
	fmt.Println(response.Body)

*/



type happyResponse struct {
	Status string
	Body string
}

type RequestParams struct {
	UrlString string
	Path string
	QueryParams []string
}

func build_endpoint_url(request RequestParams) (string) {
	// Base url. Should be the incoming param.
	apiUrl, _ := url.Parse(request.UrlString)

	// Path params
    apiUrl.Path += request.Path // path params if any

	// Query params
	params := url.Values{}
	for _,v := range request.QueryParams {
		if (strings.Contains(v, "=")) {
			newArr := strings.Split(v, "=")
			params.Add(newArr[0], newArr[1])
		}
	}
	apiUrl.RawQuery = params.Encode() 

	return apiUrl.String()
}

func Get(request RequestParams) (happyResponse) {

	apiUrl := build_endpoint_url(request)

	// Make the network call
	client := &http.Client{}
    r, _ := http.NewRequest("GET", apiUrl, nil) // URL-encoded payload
    r.Header.Add("Content-Type", "application/json")
    resp, _ := client.Do(r)

	// Read the response to return it as string
	resp_body, _ := ioutil.ReadAll(resp.Body)

	return happyResponse{resp.Status, string(resp_body)}
}

