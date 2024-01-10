package api

import (
	"encoding/csv"
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
	fmt.Printf("%s, %d, %s\n", project, days, jql)
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
	var classifications []models.Classification
	for _, issue := range issues[:10] {
		resp, err := api.openai.Classify(categories, issue.Conversation)
		if err != nil {
			return err
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

	fmt.Println("Analyzing results...")
	stats := make(map[string]map[string]int)
	for _, classification := range classifications {
		category := classification.Category
		subcategory := classification.Subcategory

		if _, ok := stats[category]; !ok {
			stats[category] = make(map[string]int)
		}

		stats[category][subcategory]++
	}
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

func generateIssuesCsv(issues []models.Classification) error {
	file, err := os.Create("issues.csv")
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	err = writer.Write([]string{"key", "summary", "category", "subcategory", "link"})
	if err != nil {
		return err
	}
	for _, issue := range issues {
		row := []string{
			issue.Key,
			issue.Summary,
			issue.Category,
			issue.Subcategory,
			fmt.Sprintf("https://jira.gustocorp.com/browse/%s", issue.Key),
		}
		fmt.Println(issue.Key,
			issue.Summary,
			issue.Category,
			issue.Subcategory,
			fmt.Sprintf("https://jira.gustocorp.com/browse/%s", issue.Key))
		err = writer.Write(row)
		if err != nil {
			return err
		}
	}
	return nil
}
func generateStatsCsv(stats map[string]map[string]int) error {
	file, err := os.Create("stats.csv")
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	err = writer.Write([]string{"category", "subcategory", "count"})
	if err != nil {
		return err
	}
	for category, subcategory := range stats {
		for sc, count := range subcategory {
			row := []string{
				category,
				sc,
				fmt.Sprintf("%d", count),
			}
			err = writer.Write(row)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
