package api

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/eddique/jcat/pkg/core/models"
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
	fmt.Println("Fetching issues...")
	issueResponse, err := api.jira.GetIssues(project, days, jql)
	if err != nil {
		return err
	}
	issues := parseIssues(issueResponse.Issues)
	fmt.Println("Parsing issues...")
	var conversations []string
	for _, issue := range issues[:10] {
		conversations = append(conversations, issue.Conversation)
	}
	fmt.Println("Creating categories...")
	samples := strings.Join(conversations, "\n********* Issue ********\n")
	categories, err := api.openai.CreateCategories(samples)
	if err != nil {
		return err
	}
	fmt.Println("Classifying issues...")
	classifications, err := api.generateClassifications(issues[:10], categories)
	if err != nil {
		return err
	}

	fmt.Println("Analyzing results...")
	stats := generateStats(classifications)

	fmt.Println("Generating Issues Report...")
	err = generateIssuesCsv(classifications)
	if err != nil {
		return err
	}
	fmt.Println("Created issues.csv!")
	fmt.Println("Generating Stats Report...")
	err = generateStatsCsv(stats)
	if err != nil {
		return err
	}
	fmt.Println("Created stats.csv!")
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

func (api ApiAdapter) generateClassifications(issues []models.IssueData, categories string) ([]models.Classification, error) {
	var classifications []models.Classification
	length := len(issues)
	statusBar := []byte(strings.Repeat("_", 25))
	for i, issue := range issues {
		fmt.Printf("Generating classification for %s %s %d of %d\n", issue.Key, string(statusBar), i+1, length)
		resp, err := api.openai.Classify(categories, issue.Conversation)
		if err != nil {
			fmt.Println(err)
			continue
		}
		var category models.Category
		err = json.Unmarshal([]byte(resp), &category)
		if err != nil {
			fmt.Println(err)
			continue
		}
		classification := models.Classification{
			Key:         issue.Key,
			Summary:     issue.Summary,
			Category:    category.Category,
			Subcategory: category.Subcategory,
		}
		classifications = append(classifications, classification)
	}
	return classifications, nil
}
