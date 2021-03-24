package brev_api

import (
	"fmt"

	"github.com/brevdev/brev-go-cli/internal/requests"
)

type ProjectPackage struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	ProjectId   string `json:"project_id"`
	CreatedDate string `json:"created_date"`
	HomePage    string `json:"home_page"`
	Status      string `json:"status"`
	Version     string `json:"version"`
}

type ProjectPackages struct {
	Packages []ProjectPackage `json:"packages"`
}

type ResponseAddPackage struct {
	Package ProjectPackage `json:"package"`
}

type ResponseRemovePackage struct {
	ID string `json:"id"`
}

func (a *Agent) GetPackages(projectID string) ([]ProjectPackage, error) {
	request := requests.RESTRequest{
		Method:   "GET",
		Endpoint: brevEndpoint("package"),
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
		return nil, fmt.Errorf("failed to get packages: %s", err)
	}

	var payload ProjectPackages
	err = response.DecodePayload(&payload)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize response payload: %s", err)
	}

	return payload.Packages, nil
}

func (a *Agent) AddPackage(projectID string, name string) (*ResponseAddPackage, error) {
	request := requests.RESTRequest{
		Method:   "POST",
		Endpoint: brevEndpoint("package"),
		QueryParams: []requests.QueryParam{
			{"utm_source", "cli"},
		},
		Headers: []requests.Header{
			{"Authorization", "Bearer " + a.Key.AccessToken},
		},
		Payload: map[string]string{
			"name":       name,
			"project_id": projectID,
		},
	}
	response, err := request.Submit()
	if err != nil {
		return nil, fmt.Errorf("failed to create package: %s", err)
	}

	var payload ResponseAddPackage
	err = response.DecodePayload(&payload)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize response payload: %s", err)
	}

	return &payload, nil
}

func (a *Agent) RemovePackage(packageID string) (*ResponseRemovePackage, error) {
	request := requests.RESTRequest{
		Method:   "DELETE",
		Endpoint: brevEndpoint("package/" + packageID),
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

	var payload ResponseRemovePackage
	err = response.DecodePayload(&payload)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}
