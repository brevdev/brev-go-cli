package brev_api

import (
	"github.com/brevdev/brev-go-cli/internal/requests"
)

type Module struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	CreateDate string `json:"create_date"`
	ProjectId  string `json:"project_id"`
	UserId     string `json:"user_id"`
	Source     string `json:"source"`
}

type Modules struct {
	Modules []Module `json:"modules"`
}

type ResponseUpdateModule struct {
	Module Module `json:"module"`
	StdOut string `json:"stdout"`
}

func (a *Agent) GetModules() (*Modules, error) {
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
	response, err := request.SubmitStrict()
	if err != nil {
		return nil, err
	}

	var payload Modules
	err = response.UnmarshalPayload(&payload)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}

func (a *Agent) UpdateModule(moduleID string, source string) (*ResponseUpdateModule, error) {
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
	response, err := request.SubmitStrict()
	if err != nil {
		return nil, err
	}

	var payload ResponseUpdateModule
	err = response.UnmarshalPayload(&payload)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}
