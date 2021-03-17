package brev

import (
	"github.com/brevdev/brev-go-cli/internal/requests"
)

type BrevProjects struct {
	Endpoints []BrevProject `json:"projects"`
}

type BrevProjectVariables struct {
	Variables []BrevProjectVariable `json:"variables"`
}

type BrevProjectPackages struct {
	Packages []BrevProjectPackage `json:"packages"`
}

type BrevProjectLogs struct {
	Logs []BrevProjectLog `json:"logs"`
}

type BrevProject struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Methods []string `json:"methods"`
	UserId  string   `json:"user_id"`
	Domain  string   `json:"domain"`
}

type BrevProjectVariable struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	ProjectId   string `json:"project_id"`
	CreatedDate string `json:"created_date"`
}

type BrevProjectPackage struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	ProjectId   string `json:"project_id"`
	CreatedDate string `json:"created_date"`
	HomePage    string `json:"home_page"`
	Status      string `json:"status"`
	Version     string `json:"version"`
}

type BrevProjectLog struct {
	LogType   string             `json:"type"`
	Timestamp string             `json:"timestamp"`
	Origin    string             `json:"origin"`
	Meta      BrevProjectLogMeta `json:"meta"`
}

type BrevProjectLogMeta struct {
	Uri           string  `json:"uri"`
	RequestId     string  `json:"request_id"`
	CpuTime       float64 `json:"cpu_time"`
	RequestMethod string  `json:"request_method"`
	StatusCode    int     `json:"status_code"`
	WallTime      float64 `json:"wall_time"`
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

func (a *BrevAgent) GetVariables(project *BrevProject) ([]BrevProjectVariable, error) {
	request := requests.RESTRequest{
		Method:   "GET",
		Endpoint: brevEndpoint("variable"),
		QueryParams: []requests.QueryParam{
			{"utm_source", "cli"},
			{"project_id", project.Id},
		},
		Headers: []requests.Header{
			{"Authorization", "Bearer " + a.Key.AccessToken},
		},
	}
	response, err := request.Submit()
	if err != nil {
		return nil, err
	}

	var payload BrevProjectVariables
	response.DecodePayload(&payload)

	return payload.Variables, nil
}

func (a *BrevAgent) GetPackages(project *BrevProject) ([]BrevProjectPackage, error) {
	request := requests.RESTRequest{
		Method:   "GET",
		Endpoint: brevEndpoint("package"),
		QueryParams: []requests.QueryParam{
			{"utm_source", "cli"},
			{"project_id", project.Id},
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

func (a *BrevAgent) GetLogs(project *BrevProject, log_type string) ([]BrevProjectLog, error) {
	request := requests.RESTRequest{
		Method:   "GET",
		Endpoint: brevEndpoint("logs"),
		QueryParams: []requests.QueryParam{
			{"utm_source", "cli"},
			{"project_id", project.Id},
			{"type", log_type},
		},
		Headers: []requests.Header{
			{"Authorization", "Bearer " + a.Key.AccessToken},
		},
	}
	response, err := request.Submit()
	if err != nil {
		return nil, err
	}

	var payload BrevProjectLogs
	response.DecodePayload(&payload)

	return payload.Logs, nil
}
