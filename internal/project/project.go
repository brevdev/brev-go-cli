package project

import (
	"fmt"

	"github.com/brevdev/brev-go-cli/internal/cmdcontext"
)

type Project struct {
	Domain     string `json:"domain"`
	CreateDate string `json:"create_date"`
	Id         string `json:"id"`
	Name       string `json:"name"`
	UserId     string `json:"user_id"`
}

func logic(context *cmdcontext.Context) error {
	fmt.Fprintln(context.Out, "project called")
	return nil
}
