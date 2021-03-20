package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type RESTRequest struct {
	Method      string
	Endpoint    string
	QueryParams []QueryParam
	Headers     []Header
	Payload     interface{}
}

type RESTResponse struct {
	StatusCode int
	Payload    []byte
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

	payloadBytes, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	return &RESTResponse{
		StatusCode: res.StatusCode,
		Payload:    payloadBytes,
	}, nil
}

// UnmarshalPayload converts the raw response body into the given interface
// Usage:
//   var foo MyStruct
//   response.UnmarshalPayload(&foo)
func (r *RESTResponse) UnmarshalPayload(v interface{}) error {
	err := json.Unmarshal(r.Payload, v)
	return err
}

// PayloadAsString returns the response body as a string
func (r *RESTResponse) PayloadAsString() (string, error) {
	return string(r.Payload), nil
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
