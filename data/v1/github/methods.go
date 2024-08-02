package github

import (
	"errors"
	"strings"
)

func (s *SetupData) Validate() error {
	if s.Repo == "" {
		return errors.New("repo is required")
	}

	if len(strings.Split(s.Repo, "/")) != 2 {
		return errors.New("invalid repo. It should be owner_name/repo_name")
	}
	return nil
}
