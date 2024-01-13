package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/eddique/jcat/pkg/core/configs"
	"github.com/eddique/jcat/pkg/core/models"
)

type JiraAdapter struct{}

func NewJiraAdapter() *JiraAdapter {
	return &JiraAdapter{}
}

func (jira JiraAdapter) FetchIssues(issues *[]models.Issue, project string, days int, jql string, startAt int, count int) error {
	fmt.Printf("Fetching page %d ...\n", count+1)
	var query string
	url := "https://jira.gustocorp.com/rest/api/2/search"
	fromDate := formatDate(days)
	if jql == "" {
		query = "project = " + project + " AND createdDate >= " + fromDate
	} else {
		query = jql
	}
	data := models.IssueQueryRequest{
		Expand:     []string{"comment"},
		Fields:     []string{"summary", "description", "comment"},
		JQL:        query,
		MaxResults: 100000,
		StartAt:    startAt,
	}
	fmt.Println(query)
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+configs.JiraApiKey())
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	issueData, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var issueResponse models.IssueQueryResponse
	err = json.Unmarshal([]byte(issueData), &issueResponse)
	if err != nil {
		return err
	}
	*issues = append(*issues, issueResponse.Issues...)
	if (issueResponse.StartAt + issueResponse.MaxResults) < issueResponse.Total {
		jira.FetchIssues(issues, project, days, jql, (issueResponse.StartAt + issueResponse.MaxResults), count+1)
	}
	return nil
}
func formatDate(days int) string {
	return time.Now().AddDate(0, 0, -days).Format("2006-01-02")
}

func parseIssues(jiraIssues []models.Issue) []models.IssueData {
	var issues []models.IssueData
	for _, jiraIssue := range jiraIssues {
		conversation := fmt.Sprintf(
			"Summary: %s \nDescription: %s\n",
			jiraIssue.Fields.Summary,
			jiraIssue.Fields.Description,
		)
		for _, comment := range jiraIssue.Fields.Comment.Comments {
			conversation += fmt.Sprintf("User: %s", comment.Body)
		}
		issue := models.IssueData{
			ID:           jiraIssue.ID,
			Key:          jiraIssue.Key,
			Summary:      jiraIssue.Fields.Summary,
			Conversation: conversation,
		}
		issues = append(issues, issue)
	}
	return issues
}
