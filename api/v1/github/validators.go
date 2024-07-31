package github

import "errors"

func (s *SetupData) Validate() error {
	if s.Owner == "" {
		return errors.New("owner is required")
	}
	if s.Repo == "" {
		return errors.New("repo is required")
	}
	return nil
}
