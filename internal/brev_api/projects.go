package brev_api

import (
	"github.com/brevdev/brev-go-cli/internal/requests"
)

type Project struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	UserId     string `json:"user_id"`
	Domain     string `json:"domain"`
	CreateDate string `json:"create_date"`
}

type Projects struct {
	Endpoints []Project `json:"projects"`
}

type ResponseCreateProject struct {
	Module  Module  `json:"module"`
	Project Project `json:"project"`
}

func (a *Agent) GetProjects() ([]Project, error) {
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
	response, err := request.SubmitStrict()
	if err != nil {
		return nil, err
	}

	var payload Projects
	err = response.UnmarshalPayload(&payload)
	if err != nil {
		return nil, err
	}

	return payload.Endpoints, nil
}

func (a *Agent) CreateProject(name string) (*ResponseCreateProject, error) {
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
	response, err := request.SubmitStrict()
	if err != nil {
		return nil, err
	}

	var payload ResponseCreateProject
	err = response.UnmarshalPayload(&payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}
