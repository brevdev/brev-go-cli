package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type RESTRequest struct {
	URI         string
	Method      string
	Endpoint    string
	QueryParams []QueryParam
	Headers     []Header
	Payload     interface{}
}

type RESTResponse struct {
	StatusCode int
	Headers    []Header
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

type RESTResponseError struct {
	RequestURI         string
	ResponseStatusCode int
}

func (e *RESTResponseError) Error() string {
	return fmt.Sprintf("REST request to `%s` returned an error status code: %d", e.RequestURI, e.ResponseStatusCode)
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
	r.URI = req.URL.String()

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

	var headers []Header
	for key, values := range res.Header {
		headers = append(headers, Header{
			Key:   key,
			Value: strings.Join(values, "\n"),
		})
	}
	return &RESTResponse{
		Headers:    headers,
		StatusCode: res.StatusCode,
		Payload:    payloadBytes,
	}, nil
}

// SubmitStrict performs the HTTP request, returning a resultant RESTResponse if the response's status code is < 400.
// Usage:
//   request = &RESTRequest{ ... }
//   response, _ := request.Submit()
func (r *RESTRequest) SubmitStrict() (*RESTResponse, error) {
	response, err := r.Submit()
	if err != nil {
		return nil, err
	}
	if response.StatusCode >= 400 {
		return nil, &RESTResponseError{
			RequestURI:         r.URI,
			ResponseStatusCode: response.StatusCode,
		}
	}
	return response, nil
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

// PayloadAsPrettyJSONString returns the response body as a formatted JSON string
// The response body must be valid JSON in the form either of a list or a map.
func (r *RESTResponse) PayloadAsPrettyJSONString() (string, error) {
	prefix := ""
	indent := "  "

	// attempt to marshal as typical JSON (e.g.: { <el>: { <el>: ... }}
	var payloadStructJson map[string]interface{}
	err := json.Unmarshal(r.Payload, &payloadStructJson)
	if err == nil {
		jsonBytes, err := json.MarshalIndent(payloadStructJson, prefix, indent)
		if err != nil {
			return "", fmt.Errorf("failed to marhsal JSON struct: %s", err)
		}
		return string(jsonBytes), nil
	}

	// error -- try to marshal again, this time as a list
	var payloadStructList []interface{}
	err = json.Unmarshal(r.Payload, &payloadStructList)
	if err == nil {
		listBytes, err := json.MarshalIndent(payloadStructList, prefix, indent)
		if err != nil {
			return "", fmt.Errorf("failed to marhsal list struct: %s", err)
		}
		return string(listBytes), nil
	}

	return "", fmt.Errorf("response was not valid JSON")
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
