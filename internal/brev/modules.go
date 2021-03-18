package brev

import (
	"github.com/brevdev/brev-go-cli/internal/requests"
)

type BrevModule struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	CreateDate string `json:"create_date"`
	ProjectId  string `json:"project_id"`
	UserId     string `json:"user_id"`
}

type BrevModules struct {
	Modules []BrevModule `json:"modules"`
}

type ResponseUpdateModule struct {
	Module BrevModule `json:"module"`
	StdOut string     `json:"stdout"`
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

func (a *BrevAgent) UpdateModule(moduleID string, source string) (*ResponseUpdateModule, error) {
	request := requests.RESTRequest{
		Method:   "PUT",
		Endpoint: brevEndpoint("module/" + moduleID),
		QueryParams: []requests.QueryParam{
			{"utm_source", "cli"},
		},
		Headers: []requests.Header{
			{"Authorization", "Bearer " + a.Key.AccessToken},
		},
		Payload: map[string]string{
			"source": source,
		},
	}
	response, err := request.Submit()
	if err != nil {
		return nil, err
	}

	var payload ResponseUpdateModule
	response.DecodePayload(&payload)

	return &payload, nil
}
