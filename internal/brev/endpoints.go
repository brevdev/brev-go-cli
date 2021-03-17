package brev

import (
	"github.com/brevdev/brev-go-cli/internal/requests"
)

type BrevEndpoints struct {
	Endpoints []BrevEndpoint `json:"endpoints"`
}

type BrevEndpoint struct {
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Methods    []string `json:"methods"`
	Uri        string   `json:"uri"`
	Archived   bool     `json:"archived"`
	CreateDate string   `json:"create_date"`
	ProjectId  string   `json:"project_id"`
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
