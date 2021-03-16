package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	Body   string
}

type RequestParams struct {
	UrlString   string
	Path        string
	QueryParams []string
}

func build_endpoint_url(request RequestParams) string {
	// Base url. Should be the incoming param.
	apiUrl, _ := url.Parse(request.UrlString)

	// Path params
	apiUrl.Path += request.Path // path params if any

	// Query params
	params := url.Values{}
	for _, v := range request.QueryParams {
		if strings.Contains(v, "=") {
			newArr := strings.Split(v, "=")
			params.Add(newArr[0], newArr[1])
		}
	}
	apiUrl.RawQuery = params.Encode()

	return apiUrl.String()
}

func Get(request RequestParams) happyResponse {

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

type PostRequestParams struct {
	RequestParams RequestParams
	Body          string
}

func Post(request PostRequestParams) happyResponse {

	apiUrl := build_endpoint_url(request.RequestParams)

	// Make the network call
	client := &http.Client{}
	r, _ := http.NewRequest("POST", apiUrl, strings.NewReader(request.Body)) // URL-encoded payload
	r.Header.Add("Content-Type", "application/json")
	resp, _ := client.Do(r)

	// Read the response to return it as string
	resp_body, _ := ioutil.ReadAll(resp.Body)

	return happyResponse{resp.Status, string(resp_body)}
}

type RESTRequest struct {
	Method      string
	Endpoint    string
	QueryParams []QueryParam
	Headers     []Header
	Payload     interface{}
}

type RESTResponse struct {
	StatusCode int
	Payload    io.ReadCloser
}

type QueryParam struct {
	Key   string
	Value string
}

type Header struct {
	Key   string
	Value string
}

// BuildHTTPRequest constructs the complete net/http object which is needed to
// perform a final HTTP request.
//
// It is not necessary to use this function, but it may be useful if the net.Request
// object needs to be inspected or modified for advanced use cases.
func (r *RESTRequest) BuildHTTPRequest() (*http.Request, error) {
	var payload io.Reader
	if r.Method == "PUT" || r.Method == "POST" || r.Method == "PATCH" {
		payloadBytes, _ := json.Marshal(r.Payload)
		payload = bytes.NewBuffer(payloadBytes)
	} else if r.Method == "GET" || r.Method == "DELETE" {
		payload = nil
	} else {
		return nil, errors.New(fmt.Sprintf("Unknown method: %s", r.Method))
	}

	// set up request
	req, err := http.NewRequest(
		r.Method,
		r.Endpoint,
		payload,
	)
	if err != nil {
		return nil, err
	}

	// build query parameters and encode
	q := req.URL.Query()
	for _, param := range r.QueryParams {
		q.Add(param.Key, param.Value)
	}
	req.URL.RawQuery = q.Encode()

	// build headers
	// TODO: remove assumed Content-Type header?
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	for _, header := range r.Headers {
		req.Header.Set(header.Key, header.Value)
	}

	return req, nil
}

// Submit performs the HTTP request, returning a resultant RESTResponse
// Usage:
//   request = &RESTRequest{ ... }
//   response, _ := request.Submit()
func (r *RESTRequest) Submit() (*RESTResponse, error) {
	req, err := r.BuildHTTPRequest()
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return &RESTResponse{
		StatusCode: res.StatusCode,
		Payload:    res.Body,
	}, nil
}

// DecodePayload converts the raw response body into the given interface
// Usage:
//   var foo MyStruct
//   response.DecodePayload(&foo)
func (r *RESTResponse) DecodePayload(v interface{}) error {
	defer r.Payload.Close()
	return json.NewDecoder(r.Payload).Decode(&v)
}

// EXAMPLE USAGE:
/*
type foo struct {
	Hi string
}

func SubmitRequestWithStruct() {
	request := &RESTRequest{
		Method:   "GET",
		Endpoint: "www.google.com",
		QueryParams: []QueryParam{
			{"foo", "bar"},
		},
		Headers: []Header{
			{"Content-Type", "application/json"},
		},
		Payload: foo{
			Hi: "there",
		},
	}
	response, _ := request.Submit()

	var myCoolResponse foo
	response.DecodePayload(&myCoolResponse)
}

func SubmitRequestWithJSON() {
	request := &RESTRequest{
		Method:   "GET",
		Endpoint: "www.google.com",
		QueryParams: []QueryParam{
			{"foo", "bar"},
		},
		Headers: []Header{
			{"Content-Type", "application/json"},
		},
		Payload: map[string]string{
			"hi": "there",
		},
	}
	response, _ := request.Submit()

	var myCoolResponse foo
	response.DecodePayload(&myCoolResponse)
}
*/
