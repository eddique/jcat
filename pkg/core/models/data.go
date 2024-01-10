package models

type IssueData struct {
	ID           string `json:"id"`
	Key          string `json:"key"`
	Conversation string `json:"conversation"`
	Summary      string `json:"summary"`
}

type Category struct {
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
}

type Classification struct {
	Key         string `json:"key"`
	Summary     string `json:"summary"`
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
}
