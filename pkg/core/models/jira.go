package models

type IssueQueryResponse struct {
	Issues  []Issue `json:"issues"`
	StartAt int     `json:"startAt"`
	Total   int     `json:"total"`
}
type Issue struct {
	ID     string      `json:"id"`
	Key    string      `json:"key"`
	Fields IssueFields `json:"fields"`
}
type IssueFields struct {
	Description string        `json:"description"`
	Summary     string        `json:"summary"`
	Comment     IssueComments `json:"comment"`
}
type IssueComments struct {
	Comments []Comment `json:"comments"`
}
type Comment struct {
	Author Author `json:"author"`
	Body   string `json:"body"`
}
type Author struct {
	Key         string `json:"key"`
	Email       string `json:"emailAddress"`
	DisplayName string `json:"displayName"`
}
type IssueQueryRequest struct {
	Expand     []string `json:"expand"`
	Fields     []string `json:"fields"`
	JQL        string   `json:"jql"`
	MaxResults int      `json:"maxResults"`
}
