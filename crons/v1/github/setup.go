package github

import (
	"log"

	"github.com/robfig/cron/v3"
)

var Cron = cron.New()

func AddJob(schedule string, job func()) {
	/*
		This function adds a job to the cron job runner
	*/
	log.Println("add new job")
	_, err := Cron.AddFunc(schedule, job)
	if err != nil {
		log.Printf("Error adding cron job: %v", err)
	}
}

func StartCron() {
	/*
		This function starts the cron job runner.
	*/
	go func() {
		Cron.Start()
		log.Println("Cron jobs started")
	}()

}
