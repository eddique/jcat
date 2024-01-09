package ports

import "github.com/eddique/jcat/pkg/core/models"

type IssuePort interface {
	GetIssues(project string, days int, jql string) (*models.IssueQueryResponse, error)
}
