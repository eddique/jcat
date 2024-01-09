package main

import (
	"fmt"

	"github.com/eddique/jcat/pkg/api"
)

func main() {
	jira := api.NewJiraAdapter()
	// issues, err := jira.GetIssues("IT", 90, "")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// for _, issue := range issues.Issues {
	// 	fmt.Printf("\nKey: %s,\nSummary: %s\n\n", issue.Key, issue.Fields.Summary)
	// }
	// fmt.Println("Done!")
	openai := api.NewOpenAIAdapter()
	app := api.NewApiAdapter(openai, jira)
	err := app.Run()
	if err != nil {
		fmt.Println(err)
	}
}
