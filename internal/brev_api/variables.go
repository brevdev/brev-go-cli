package brev_api

import (
	"github.com/brevdev/brev-go-cli/internal/requests"
)

type ProjectVariable struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	ProjectId  string `json:"project_id"`
	CreateDate string `json:"create_date"`
}

type ProjectVariables struct {
	Variables []ProjectVariable `json:"variables"`
}

type ResponseAddVariable struct {
	Variable ProjectVariable `json:"variable"`
}

type ResponseRemoveVariable struct {
	ID string `json:"id"`
}

func (a *Agent) GetVariables(projectID string) ([]ProjectVariable, error) {
	request := requests.RESTRequest{
		Method:   "GET",
		Endpoint: brevEndpoint("variable"),
		QueryParams: []requests.QueryParam{
			{"utm_source", "cli"},
			{"project_id", projectID},
		},
		Headers: []requests.Header{
			{"Authorization", "Bearer " + a.Key.AccessToken},
		},
	}
	response, err := request.Submit()
	if err != nil {
		return nil, err
	}

	var payload ProjectVariables
	err = response.UnmarshalPayload(&payload)
	if err != nil {
		return nil, err
	}

	return payload.Variables, nil
}

func (a *Agent) AddVariable(projectID string, name string, value string) (*ResponseAddVariable, error) {
	request := requests.RESTRequest{
		Method:   "POST",
		Endpoint: brevEndpoint("variable"),
		QueryParams: []requests.QueryParam{
			{"utm_source", "cli"},
		},
		Headers: []requests.Header{
			{"Authorization", "Bearer " + a.Key.AccessToken},
		},
		Payload: map[string]string{
			"name":       name,
			"value":      value,
			"project_id": projectID,
		},
	}
	response, err := request.Submit()
	if err != nil {
		return nil, err
	}

	var payload ResponseAddVariable
	err = response.UnmarshalPayload(&payload)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}

func (a *Agent) RemoveVariable(variableID string) (*ResponseRemoveVariable, error) {
	request := requests.RESTRequest{
		Method:   "DELETE",
		Endpoint: brevEndpoint("variable/" + variableID),
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

	var payload ResponseRemoveVariable
	err = response.UnmarshalPayload(&payload)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}
