package brev_api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/brevdev/brev-go-cli/internal/requests"
	"github.com/gorilla/websocket"
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

func (a *Agent) GetHistoricalLogs(projectID string, logType string) ([]ProjectLog, error) {
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
	response, err := request.SubmitStrict()
	if err != nil {
		return nil, err
	}

	var payload ProjectLogs
	err = response.UnmarshalPayload(&payload)
	if err != nil {
		return nil, err
	}

	return payload.Logs, nil
}

func (a *Agent) TailLiveLogs(instanceId string, task string) <-chan string {

	url := brevLogEndpoint(fmt.Sprintf("?client_type=%s&instance_id=%s&task=%s", "USER", instanceId, task))

	log.Printf("connecting to %s", url)

	c, resp, err := websocket.DefaultDialer.Dial(url, http.Header{"token": []string{a.Key.AccessToken}})
	if err != nil {
		log.Printf("handshake failed with status %d", resp.StatusCode)
		log.Fatal("dial: ", err)
	}

	logs := make(chan string)

	go func() {
		defer close(logs)
		defer c.Close()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			logs <- string(message)
		}
	}()
	return logs

}
