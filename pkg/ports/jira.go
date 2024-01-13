package ports

import "github.com/eddique/jcat/pkg/core/models"

type IssuePort interface {
	FetchIssues(issues *[]models.Issue, project string, days int, jql string, startAt int, count int) error
}
