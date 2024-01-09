package ports

type GPTPort interface {
	CreateCategories(samples string) (*[]string, error)
	Classify(categories interface{}, issue string) (string, error)
}
