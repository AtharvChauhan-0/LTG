package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func worker(id int, url string, requests int, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{}

	for i := 0; i < requests; i++ {

		start := time.Now()
		resp, err := client.Get(url)

		entry := LogEntry{
			Timestamp: time.Now(),
			VU:        id,
			URL:       url,
			Method:    "GET",
		}

		if err != nil {
			entry.Error = err.Error()
		} else {
			entry.Status = resp.StatusCode
			entry.DurationMS = time.Since(start).Milliseconds()
			resp.Body.Close()
		}

		// Write log entry
		WriteLog(entry)
	}
}

func main() {
	url := "https://httpbin.org/get"
	virtualUsers := 20
	requestsPerUser := 10

	// Init logger
	err := InitLogger()
	if err != nil {
		panic(err)
	}
	defer CloseLogger()

	var wg sync.WaitGroup

	fmt.Printf("Starting load test: %d users × %d requests\n", virtualUsers, requestsPerUser)

	for i := 0; i < virtualUsers; i++ {
		wg.Add(1)
		go worker(i, url, requestsPerUser, &wg)
	}

	wg.Wait()

	fmt.Println("\n=== Test Finished ===")
	fmt.Println("Logs saved to: requests.log")
}
