package api

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/eddique/jcat/pkg/core/models"
	"github.com/eddique/jcat/pkg/ports"
)

type ApiAdapter struct {
	openai ports.GPTPort
	jira   ports.IssuePort
	rate   *rate.Limiter
}

func NewApiAdapter(openai ports.GPTPort, jira ports.IssuePort) *ApiAdapter {
	rateLimit := rate.NewLimiter(rate.Every(time.Minute), 1000)
	return &ApiAdapter{openai, jira, rateLimit}
}

func (api ApiAdapter) ClassifyIssues(project string, days int, jql string) error {
	fmt.Println("Fetching issues...")
	var issueData []models.Issue
	err := api.jira.FetchIssues(&issueData, project, days, jql, 0, 0)
	if err != nil {
		return err
	}
	issues := parseIssues(issueData)
	fmt.Printf("Fetched %d issues...\n", len(issues))
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
	var project string

	flag.StringVar(&jql, "jql", "", "Specify custom jql query")
	flag.IntVar(&days, "days", 90, "Days before now to query issues")
	flag.StringVar(&project, "project", "IT", "Specify the project")
	flag.Parse()

	err := api.ClassifyIssues(project, days, jql)
	if err != nil {
		return err
	}
	return nil
}
func (api ApiAdapter) GenerateClassifications(issues []models.IssueData, categories string) ([]models.Classification, error) {
	var classifications []models.Classification
	length := len(issues)
	statusBar := []byte(strings.Repeat("_", 25))
	for i, issue := range issues {
		fmt.Printf("Generating classification for %s %s %d of %d\n", issue.Key, string(statusBar), i+1, length)
		category, err := api.openai.Classify(categories, issue.Conversation)
		if err != nil {
			fmt.Println(err)
			continue
		}
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
func (api ApiAdapter) generateClassifications(issues []models.IssueData, categories string) ([]models.Classification, error) {
	const numWorkers = 20
	issueChan := make(chan models.IssueData, len(issues))
	classificationChan := make(chan models.Classification, len(issues))
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for issue := range issueChan {
				if err := api.rate.Wait(context.Background()); err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("Classifying issue", issue.Key)
				category, err := api.openai.Classify(categories, issue.Conversation)
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
				classificationChan <- classification
				time.Sleep(time.Second)
			}
		}()
	}
	for _, issue := range issues {
		issueChan <- issue
	}
	close(issueChan)

	go func() {
		wg.Wait()
		close(classificationChan)
	}()

	var classifications []models.Classification
	for classification := range classificationChan {
		classifications = append(classifications, classification)
	}

	return classifications, nil
}
