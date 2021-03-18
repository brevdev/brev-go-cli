package brev

import (
	"github.com/brevdev/brev-go-cli/internal/requests"
)

type BrevEndpoint struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Methods    []string `json:"methods"`
	Uri        string   `json:"uri"`
	Archived   bool     `json:"archived"`
	CreateDate string   `json:"create_date"`
	ProjectId  string   `json:"project_id"`
}

type BrevEndpoints struct {
	Endpoints []BrevEndpoint `json:"endpoints"`
}

type RequestUpdateEndpoint struct {
	Name    string   `json:"name"`
	Methods []string `json:"methods"`
	Code    string   `json:"code"`
}

type ResponseUpdateEndpoint struct {
	Endpoint BrevEndpoint `json:"endpoint"`
}

type ResponseRemoveEndpoint struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
}

func (a *BrevAgent) GetEndpoints() (*BrevEndpoints, error) {
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

	var payload BrevEndpoints
	response.DecodePayload(&payload)

	return &payload, nil
}

func (a *BrevAgent) UpdateEndpoint(endpointID string, updateRequest RequestUpdateEndpoint) (*ResponseUpdateEndpoint, error) {
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
	response.DecodePayload(&payload)

	return &payload, nil
}

func (a *BrevAgent) RemoveEndpoint(endpointID string) (*ResponseRemoveEndpoint, error) {
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
	response.DecodePayload(&payload)

	return &payload, nil
}
