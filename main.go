package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	Duration time.Duration
	Error    error
}

func worker(id int, url string, requests int, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{}

	for i := 0; i < requests; i++ {
		start := time.Now()

		resp, err := client.Get(url)
		if err != nil {
			results <- Result{Duration: 0, Error: err}
			continue
		}

		resp.Body.Close()

		results <- Result{Duration: time.Since(start), Error: nil}
	}
}

func main() {
	url := "https://httpbin.org/get"
	virtualUsers := 20
	requestsPerUser := 10

	results := make(chan Result, virtualUsers*requestsPerUser)
	var wg sync.WaitGroup

	startTest := time.Now()

	fmt.Printf("Starting load test: %d users × %d requests\n", virtualUsers, requestsPerUser)

	for i := 0; i < virtualUsers; i++ {
		wg.Add(1)
		go worker(i, url, requestsPerUser, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var totalRequests, success, failed int
	var totalDuration time.Duration

	for res := range results {
		totalRequests++

		if res.Error != nil {
			failed++
		} else {
			success++
			totalDuration += res.Duration
		}
	}

	testDuration := time.Since(startTest)

	fmt.Println("\n=== Load Test Results ===")
	fmt.Printf("Total Requests: %d\n", totalRequests)
	fmt.Printf("Success: %d\n", success)
	fmt.Printf("Failed: %d\n", failed)
	fmt.Printf("Total Test Duration: %v\n", testDuration)

	if success > 0 {
		fmt.Printf("Avg Response Time: %v\n", totalDuration/time.Duration(success))
		fmt.Printf("Requests per second: %.2f\n", float64(totalRequests)/testDuration.Seconds())
	}

	fmt.Println("=========================")
}
