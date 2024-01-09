package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func OpenAIApiKey() string {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
	return os.Getenv("OPENAI_API_KEY")
}

func JiraApiKey() string {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
	return os.Getenv("JIRA_API_KEY")
}
