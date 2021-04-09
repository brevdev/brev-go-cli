package project

import (
	"github.com/brevdev/brev-go-cli/internal/terminal"
)

type Project struct {
	Domain     string `json:"domain"`
	CreateDate string `json:"create_date"`
	Id         string `json:"id"`
	Name       string `json:"name"`
	UserId     string `json:"user_id"`
}

func logic(t *terminal.Terminal) error {
	t.Vprint("project called")
	return nil
}
