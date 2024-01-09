package ports

type GPTPort interface {
	CreateCategories(number int) (*[]string, error)
	Classify(categories []string, text string) (string, error)
}
