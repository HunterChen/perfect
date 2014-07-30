package perfect

type GithubIssue struct {
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Labels []string `json:"labels,omitempty"`
}

type GithubResponse struct {
	Url         string `json:"url"`
	LabelsUrl   string `json:"labels_url"`
	CommentsUrl string `json:"comments_url"`
	HtmlUrl     string `json:"html_url"`
	IssueNumber int    `json:"number"`
}
