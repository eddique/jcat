package api

import (
	"flag"
	"fmt"
	"os"

	"github.com/eddique/jcat/pkg/ports"
)

type ApiAdapter struct {
	openai ports.GPTPort
	jira   ports.IssuePort
}

func NewApiAdapter(openai ports.GPTPort, jira ports.IssuePort) *ApiAdapter {
	return &ApiAdapter{openai, jira}
}

func (api ApiAdapter) ClassifyIssues(project string, days int, jql string) error {
	fmt.Printf("%s, %d, %s", project, days, jql)
	issueResponse, err := api.jira.GetIssues(project, days, jql)
	if err != nil {
		return err
	}
	for _, issue := range issueResponse.Issues {
		fmt.Printf("\nKey: %s,\nSummary: %s\n\n", issue.Key, issue.Fields.Summary)
	}
	fmt.Println("Done!")

	return nil
}

func (api ApiAdapter) Run() error {
	var jql string
	var days int

	flag.StringVar(&jql, "jql", "", "Specify custom jql query")
	flag.IntVar(&days, "days", 90, "Days before now to query issues")
	flag.Parse()

	project := os.Args[1]
	err := api.ClassifyIssues(project, days, jql)
	if err != nil {
		return err
	}
	return nil
}
