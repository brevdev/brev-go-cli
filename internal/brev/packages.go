package brev

import (
	"github.com/brevdev/brev-go-cli/internal/requests"
)

type BrevProjectPackage struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	ProjectId   string `json:"project_id"`
	CreatedDate string `json:"created_date"`
	HomePage    string `json:"home_page"`
	Status      string `json:"status"`
	Version     string `json:"version"`
}

type BrevProjectPackages struct {
	Packages []BrevProjectPackage `json:"packages"`
}

type ResponseAddPackage struct {
	Package BrevProjectPackage `json:"package"`
}

type ResponseRemovePackage struct {
	ID string `json:"id"`
}

func (a *BrevAgent) GetPackages(projectID string) ([]BrevProjectPackage, error) {
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
		return nil, err
	}

	var payload BrevProjectPackages
	response.DecodePayload(&payload)

	return payload.Packages, nil
}

func (a *BrevAgent) AddPackage(projectID string, name string) (*ResponseAddPackage, error) {
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
		return nil, err
	}

	var payload ResponseAddPackage
	response.DecodePayload(&payload)

	return &payload, nil
}

func (a *BrevAgent) RemovePackage(packageID string) (*ResponseRemovePackage, error) {
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
	response.DecodePayload(&payload)

	return &payload, nil
}
