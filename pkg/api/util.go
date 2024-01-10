package api

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/eddique/jcat/pkg/core/models"
)

func generateStats(classifications []models.Classification) map[string]map[string]int {
	stats := make(map[string]map[string]int)
	for _, classification := range classifications {
		category := classification.Category
		subcategory := classification.Subcategory

		if _, ok := stats[category]; !ok {
			stats[category] = make(map[string]int)
		}

		stats[category][subcategory]++
	}
	return stats
}
func generateIssuesCsv(issues []models.Classification) error {
	file, err := os.Create("issues.csv")
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	err = writer.Write([]string{"key", "summary", "category", "subcategory", "link"})
	if err != nil {
		return err
	}
	for _, issue := range issues {
		row := []string{
			issue.Key,
			issue.Summary,
			issue.Category,
			issue.Subcategory,
			fmt.Sprintf("https://jira.gustocorp.com/browse/%s", issue.Key),
		}
		err = writer.Write(row)
		if err != nil {
			return err
		}
	}
	return nil
}
func generateStatsCsv(stats map[string]map[string]int) error {
	file, err := os.Create("stats.csv")
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	err = writer.Write([]string{"category", "subcategory", "count"})
	if err != nil {
		return err
	}
	for category, subcategory := range stats {
		for sc, count := range subcategory {
			row := []string{
				category,
				sc,
				fmt.Sprintf("%d", count),
			}
			err = writer.Write(row)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
