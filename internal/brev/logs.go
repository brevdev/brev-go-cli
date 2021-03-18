package brev

import "github.com/brevdev/brev-go-cli/internal/requests"

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

type BrevProjectLogs struct {
	Logs []BrevProjectLog `json:"logs"`
}

func (a *BrevAgent) GetLogs(projectID string, logType string) ([]BrevProjectLog, error) {
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

	var payload BrevProjectLogs
	response.DecodePayload(&payload)

	return payload.Logs, nil
}
