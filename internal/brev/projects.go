package brev

import (
	"github.com/brevdev/brev-go-cli/internal/requests"
)

type BrevProject struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	UserId     string `json:"user_id"`
	Domain     string `json:"domain"`
	CreateDate string `json:"create_date"`
}

type BrevProjects struct {
	Endpoints []BrevProject `json:"projects"`
}

type ResponseCreateProject struct {
	Module  BrevModule  `json:"module"`
	Project BrevProject `json:"project"`
}

func (a *BrevAgent) GetProjects() ([]BrevProject, error) {
	request := requests.RESTRequest{
		Method:   "GET",
		Endpoint: brevEndpoint("_project"),
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

	var payload BrevProjects
	response.DecodePayload(&payload)

	return payload.Endpoints, nil
}

func (a *BrevAgent) CreateProject(name string) (*ResponseCreateProject, error) {
	request := requests.RESTRequest{
		Method:   "POST",
		Endpoint: brevEndpoint("_project"),
		QueryParams: []requests.QueryParam{
			{"utm_source", "cli"},
		},
		Headers: []requests.Header{
			{"Authorization", "Bearer " + a.Key.AccessToken},
		},
		Payload: map[string]string{
			"name": name,
		},
	}
	response, err := request.Submit()
	if err != nil {
		return nil, err
	}

	var payload ResponseCreateProject
	response.DecodePayload(&payload)
	return &payload, nil
}
