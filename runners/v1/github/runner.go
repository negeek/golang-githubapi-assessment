package github

import (
	"log"
)

func Run(workerCount int, funcs ...func()) {
	/*
		This function runs the functions passed to it using a worker pool pattern.
		It doesn't return the result of the functions; it is meant for functions that just run without returning any value.


		Arguments:
			workerCount:    number of workers to limit the concurrent goroutines
			funcs: a list of functions to be run immediately this function is called
	*/
	if workerCount <= 0 {
		workerCount = 1
	}
	if workerCount > 10 {
		workerCount = 10
	}
	if len(funcs) <= 0 {
		return
	}

	jobs := make(chan func(), len(funcs))
	results := make(chan struct{}, len(funcs))

	// Start the worker pool
	for w := 0; w < workerCount; w++ {
		go worker(jobs, results)
	}

	// Send the functions to the job channel
	for _, fn := range funcs {
		jobs <- fn
	}
	close(jobs)

	// Wait for all results to be processed
	for a := 0; a < len(funcs); a++ {
		<-results
	}

}

func worker(jobs <-chan func(), results chan<- struct{}) {
	for job := range jobs {
		log.Println("running function")
		job()
		results <- struct{}{}
	}
}
