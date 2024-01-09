package api

type OpenAIAdapter struct{}

func NewOpenAIAdapter() *OpenAIAdapter {
	return &OpenAIAdapter{}
}

func (gpt OpenAIAdapter) Classify(categories []string, text string) (string, error) {
	return "", nil
}

func (gpt OpenAIAdapter) CreateCategories(number int) (*[]string, error) {
	return nil, nil
}
