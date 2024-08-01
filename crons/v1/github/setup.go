package github

import (
	"log"

	"github.com/robfig/cron/v3"
)

var Cron = cron.New()

func AddFunc(schedule string, job func()) {
	log.Println("add new job")
	_, err := Cron.AddFunc(schedule, job)
	if err != nil {
		log.Printf("Error adding cron job: %v", err)
	}
}

func Start(immediateFuncs ...func()) {
	// Start the cron scheduler in a separate goroutine
	go func() {
		Cron.Start()
		log.Println("Cron jobs started")
	}()

	// Run immediate functions
	if len(immediateFuncs) > 0 {
		for _, fn := range immediateFuncs {
			log.Println("Running immediate function")
			fn()
		}
	}
}
