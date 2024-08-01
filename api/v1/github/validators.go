package github

import "errors"

func (t *TopNCommitAuthorsRequestData) Validate() error {
	if t.Repo == "" {
		return errors.New("repo is required")
	}
	if t.TopN == 0 {
		return errors.New("top_n is required")
	}
	return nil
}

func (r *RepoCommitsRequestData) Validate() error {
	if r.Repo == "" {
		return errors.New("repo is required")
	}
	return nil
}
