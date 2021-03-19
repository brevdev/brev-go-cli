package brev

import (
	"github.com/brevdev/brev-go-cli/internal/requests"
)

type Endpoint struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Methods    []string `json:"methods"`
	Uri        string   `json:"uri"`
	Archived   bool     `json:"archived"`
	CreateDate string   `json:"create_date"`
	ProjectId  string   `json:"project_id"`
	Code       string   `json:"code"`
}

type Endpoints struct {
	Endpoints []Endpoint `json:"endpoints"`
}

type RequestCreateEndpoint struct {
	Name      string   `json:"name"`
	Methods   []string `json:"methods"`
	Code      string   `json:"code"`
	ProjectId string   `json:"project_id"`
	Uri       string   `json:"uri"`
}

type RequestUpdateEndpoint struct {
	Name    string   `json:"name"`
	Methods []string `json:"methods"`
	Code    string   `json:"code"`
}

type ResponseUpdateEndpoint struct {
	Endpoint Endpoint `json:"endpoint"`
}

type ResponseRemoveEndpoint struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
}

func (a *Agent) GetEndpoints() (*Endpoints, error) {
	request := requests.RESTRequest{
		Method:   "GET",
		Endpoint: brevEndpoint("_endpoint"),
		QueryParams: []requests.QueryParam{
			{"utm_source", "cli"},
		},
		Headers: []requests.Header{
			{"Authorization", "Bearer " + a.Key.AccessToken},
		},
	}
	response, err := request.Submit()
	if err != nil {
		return nil, err
	}

	var payload Endpoints
	err = response.DecodePayload(&payload)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}

const dummyCode = `import variables
import shared
import variables
from global_storage import storage_context

def get():
    return {"response": "hi get"}

def post():
    return {"response": "hi post"}

`

func (a *Agent) CreateEndpoint(name string, projectId string) (*ResponseUpdateEndpoint, error) {
	request := &requests.RESTRequest{
		Method:   "POST",
		Endpoint: brevEndpoint("_endpoint"),
		QueryParams: []requests.QueryParam{
			{"utm_source", "cli"},
		},
		Headers: []requests.Header{
			{"Authorization", "Bearer " + a.Key.AccessToken},
		},
		Payload: RequestCreateEndpoint{
			Name:      name,
			ProjectId: projectId,
			Methods:   []string{},
			Code:      dummyCode,
			Uri:       "/" + name,
		},
	}

	response, err := request.Submit()
	if err != nil {
		return nil, err
	}

	var payload ResponseUpdateEndpoint
	err = response.DecodePayload(&payload)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}

func (a *Agent) UpdateEndpoint(endpointID string, updateRequest RequestUpdateEndpoint) (*ResponseUpdateEndpoint, error) {
	request := requests.RESTRequest{
		Method:   "PUT",
		Endpoint: brevEndpoint("_endpoint/" + endpointID),
		QueryParams: []requests.QueryParam{
			{"utm_source", "cli"},
		},
		Headers: []requests.Header{
			{"Authorization", "Bearer " + a.Key.AccessToken},
		},
		Payload: updateRequest,
	}
	response, err := request.Submit()
	if err != nil {
		return nil, err
	}

	var payload ResponseUpdateEndpoint
	err = response.DecodePayload(&payload)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}

func (a *Agent) RemoveEndpoint(endpointID string) (*ResponseRemoveEndpoint, error) {
	request := requests.RESTRequest{
		Method:   "DELETE",
		Endpoint: brevEndpoint("_endpoint/" + endpointID),
		QueryParams: []requests.QueryParam{
			{"utm_source", "cli"},
		},
		Headers: []requests.Header{
			{"Authorization", "Bearer " + a.Key.AccessToken},
		},
	}
	response, err := request.Submit()
	if err != nil {
		return nil, err
	}

	var payload ResponseRemoveEndpoint
	err = response.DecodePayload(&payload)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}
