package project

import "fmt"

type Project struct {
	Domain     string `json:"domain"`
	CreateDate string `json:"create_date"`
	Id         string `json:"id"`
	Name       string `json:"name"`
	UserId     string `json:"user_id"`
}

func logic() {
	fmt.Println("project called")
}
