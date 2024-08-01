package github

import (
	"log"

	"github.com/robfig/cron/v3"
)

var Cron = cron.New()

func AddFunc(schedule string, job func()) {
	_, err := Cron.AddFunc(schedule, job)
	if err != nil {
		log.Printf("Error adding cron job: %v", err)
	}
}

func Start() {
	Cron.Start()
}
