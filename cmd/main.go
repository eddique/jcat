package main

import (
	"fmt"

	"github.com/eddique/jcat/pkg/api"
)

func main() {
	jira := api.NewJiraAdapter()
	openai := api.NewOpenAIAdapter()
	app := api.NewApiAdapter(openai, jira)
	err := app.Run()
	if err != nil {
		fmt.Println(err)
	}
}
