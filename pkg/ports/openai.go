package ports

type GPTPort interface {
	CreateCategories(samples string) (string, error)
	Classify(categories string, issue string) (string, error)
}
