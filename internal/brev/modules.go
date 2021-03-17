package brev

import (
	"github.com/brevdev/brev-go-cli/internal/requests"
)

type BrevModules struct {
	Modules []BrevModule `json:"modules"`
}

type BrevModule struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	CreateDate string `json:"create_date"`
	ProjectId  string `json:"project_id"`
}

func (a *BrevAgent) GetModules() (*BrevModules, error) {
	request := requests.RESTRequest{
		Method:   "GET",
		Endpoint: brevEndpoint("module"),
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

	var payload BrevModules
	response.DecodePayload(&payload)

	return &payload, nil
}
