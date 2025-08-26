package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bernardoazevedo/faturas/internal/message"
	"github.com/joho/godotenv"
)

type Job struct {
	ID      int
	Message string
}

type Result struct {
	JobID  int
	Status bool
}

func main() {
	log.SetPrefix("workers/message: ")
	log.SetFlags(0)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = message.Start()
	if err != nil {
		log.Fatal("Error connecting to rabbitmq: " + err.Error())
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// var wg sync.WaitGroup

	// wg.Add(1)
	// amqpMessages, err := message.RetornaDelivery("notifications")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// go func() {
	// 	// defer wg.Done()
	// 	// go func() {
	// 	for message := range amqpMessages {
	// 		message := fmt.Sprintf("send: %v\n", string(message.Body))
	// 		log.Printf("write: %s", message)
	// 	}
	// 	// }()
	// }()

	// log.Println("[*] Monitoring messages. Press CTRL+C to exit")
	// <-sigchan

	// // wg.Wait()

	// defer message.AMQPconn.Close()

	// log.Println("Killed, shutting down")

	const jobCount = 100    // Total number of jobs to process
	const workerCount = 3  // Number of workers to process the jobs

	fmt.Println("Starting batch processing with synchronized result collection...")
	dispatcher(jobCount, workerCount)
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		results <- Result{JobID: job.ID, Status: true}
	}
}

func collectResults(results <-chan Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for result := range results {
		fmt.Printf("Job ID: %d, Input: %d, result: %t\n", result.JobID, result.JobID, result.Status)
	}
}

func dispatcher(jobCount, workerCount int) {
	jobs := make(chan Job, jobCount)
	results := make(chan Result, jobCount)

	var wg sync.WaitGroup

	// Start workers
	wg.Add(workerCount)
	for w := 1; w <= workerCount; w++ {
	go worker(w, jobs, results, &wg)
	}

	// Start collecting results
	var resultsWg sync.WaitGroup
	resultsWg.Add(1)
	go collectResults(results, &resultsWg)

	// Distribute jobs and wait for completion
	for j := 1; j <= jobCount; j++ {
	jobs <- Job{ID: j, Message: "aaa"}
	}
	close(jobs)
	wg.Wait()
	close(results)

	// Ensure all results are collected
	resultsWg.Wait()
}