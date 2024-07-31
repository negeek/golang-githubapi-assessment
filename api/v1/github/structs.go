package github

import "time"

type SetupData struct {
	Owner    string    `json:"owner"`
	Repo     string    `json:"repo"`
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
}
