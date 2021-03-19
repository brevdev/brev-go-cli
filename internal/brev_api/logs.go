package brev_api

import (
	"github.com/brevdev/brev-go-cli/internal/requests"
)

type ProjectLog struct {
	LogType   string         `json:"type"`
	Timestamp string         `json:"timestamp"`
	Origin    string         `json:"origin"`
	Meta      ProjectLogMeta `json:"meta"`
}

type ProjectLogMeta struct {
	Uri           string  `json:"uri"`
	RequestId     string  `json:"request_id"`
	CpuTime       float64 `json:"cpu_time"`
	RequestMethod string  `json:"request_method"`
	StatusCode    int     `json:"status_code"`
	WallTime      float64 `json:"wall_time"`
}

type ProjectLogs struct {
	Logs []ProjectLog `json:"logs"`
}

func (a *Agent) GetLogs(projectID string, logType string) ([]ProjectLog, error) {
	request := requests.RESTRequest{
		Method:   "GET",
		Endpoint: brevEndpoint("logs"),
		QueryParams: []requests.QueryParam{
			{"utm_source", "cli"},
			{"project_id", projectID},
			{"type", logType},
		},
		Headers: []requests.Header{
			{"Authorization", "Bearer " + a.Key.AccessToken},
		},
	}
	response, err := request.Submit()
	if err != nil {
		return nil, err
	}

	var payload ProjectLogs
	err = response.DecodePayload(&payload)
	if err != nil {
		return nil, err
	}

	return payload.Logs, nil
}
