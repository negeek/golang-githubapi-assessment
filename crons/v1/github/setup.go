package github

import (
	"log"

	"github.com/robfig/cron/v3"
)

var Cron = cron.New()

func AddFunc(schedule string, job func()) {
	/*
		This function adds a job to the cron job runner
	*/
	log.Println("add new job")
	_, err := Cron.AddFunc(schedule, job)
	if err != nil {
		log.Printf("Error adding cron job: %v", err)
	}
}

func Start(immediateFuncs ...func()) {
	/*
		This function starts the cron job runner.

		Arguement:
			immediateFuncs:  a list of functions to be run immediately this function is called
							essence is to avoid waiting for cron job runner schedule time.
	*/
	go func() {
		Cron.Start()
		log.Println("Cron jobs started")
	}()

	if len(immediateFuncs) > 0 {
		for _, fn := range immediateFuncs {
			log.Println("Running immediate function")
			fn()
		}
	}
}
