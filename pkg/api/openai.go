package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eddique/jcat/pkg/core/configs"
	"github.com/eddique/jcat/pkg/core/models"
	openai "github.com/sashabaranov/go-openai"
)

type OpenAIAdapter struct {
	client *openai.Client
}

func NewOpenAIAdapter() *OpenAIAdapter {
	client := openai.NewClient(configs.OpenAIApiKey())
	return &OpenAIAdapter{client}
}

func (gpt OpenAIAdapter) Classify(categories string, issue string) (*models.Category, error) {
	prompt := fmt.Sprintf(`
	Use the following categories to classify the following JIRA issue:
	%s
	return a JSON object with the keys "category" and "subcategory", 
	importantly, do not explain why just return one JSON object, without backticks, and make sure it's one of the
	categories/subcategories.

	%s`, string(categories), issue)
	resp, err := gpt.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT432K,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}
	var category models.Category
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &category)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (gpt OpenAIAdapter) CreateCategories(samples string) (string, error) {
	prompt := fmt.Sprintf(`
	The following are random samples from a JIRA project. Please create
	classifications for the category these issues could fall into, and
	sub categories underneath. These should be general enough that you 
	can classify any issue from this project into one of these, but 
	should make sense. Return a JSON object with a key "categories" and a list
	of objects with the key "category", and "subcategories" which would be a string
	list of the subcategories. Make sure to include other as a category and subcategory for each.
	
	%s`, samples)
	resp, err := gpt.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT432K,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}
