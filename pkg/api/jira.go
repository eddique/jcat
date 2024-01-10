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

func (jira JiraAdapter) GetIssues(project string, days int, jql string) (*models.IssueQueryResponse, error) {
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
		MaxResults: 10,
	}
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+configs.JiraApiKey())
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	issueData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var issueResponse models.IssueQueryResponse
	err = json.Unmarshal([]byte(issueData), &issueResponse)
	if err != nil {
		return nil, err
	}
	return &issueResponse, nil
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
