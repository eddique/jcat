package ports

import "github.com/eddique/jcat/pkg/core/models"

type GPTPort interface {
	CreateCategories(samples string) (string, error)
	Classify(categories string, issue string) (*models.Category, error)
}
