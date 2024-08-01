package github

type TopNCommitAuthorsRequestData struct {
	Repo string `json:"repo"`
	TopN int    `json:"top_n"`
}

type RepoCommitsRequestData struct {
	Repo string `json:"repo"`
}
