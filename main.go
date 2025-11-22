package main

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var (
	successCount int64
	errorCount   int64
)

func worker(id int, url string, requests int, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	for i := 0; i < requests; i++ {
		start := time.Now()
		resp, err := client.Get(url)

		entry := LogEntry{
			VU:     id,
			URL:    url,
			Method: "GET",
			Level:  "info",
		}

		if err != nil {
			entry.Error = err.Error()
			entry.Level = "error"
			atomic.AddInt64(&errorCount, 1)
		} else {
			entry.Status = resp.StatusCode
			entry.DurationMS = time.Since(start).Milliseconds()

			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				atomic.AddInt64(&successCount, 1)
			} else {
				entry.Level = "warning"
				atomic.AddInt64(&errorCount, 1)
			}

			resp.Body.Close()
		}

		// Write log entry
		WriteLog(entry)
	}
}

func main() {
	// Configuration
	url := "https://httpbin.org/get" // Changed to working endpoint
	virtualUsers := 20
	requestsPerUser := 10
	totalRequests := virtualUsers * requestsPerUser

	// Init logger
	filename, err := InitLogger()
	if err != nil {
		panic(err)
	}
	defer CloseLogger()

	fmt.Printf("Starting load test\n")
	fmt.Printf("Target URL: %s\n", url)
	fmt.Printf("Virtual Users: %d\n", virtualUsers)
	fmt.Printf("Requests per User: %d\n", requestsPerUser)
	fmt.Printf("Total Requests: %d\n", totalRequests)
	fmt.Printf("Log file: %s\n\n", filename)

	// Log test start metadata
	WriteMetadata(TestMetadata{
		Level:           "info",
		Message:         "Load test started",
		VirtualUsers:    virtualUsers,
		RequestsPerUser: requestsPerUser,
		TotalRequests:   totalRequests,
		TargetURL:       url,
	})

	testStart := time.Now()
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < virtualUsers; i++ {
		wg.Add(1)
		go worker(i, url, requestsPerUser, &wg)
	}

	wg.Wait()
	testDuration := time.Since(testStart)

	// Calculate statistics
	success := atomic.LoadInt64(&successCount)
	errors := atomic.LoadInt64(&errorCount)
	successRate := float64(success) / float64(totalRequests) * 100
	reqPerSec := float64(totalRequests) / testDuration.Seconds()

	// Log test completion metadata
	WriteMetadata(TestMetadata{
		Level: "info",
		Message: fmt.Sprintf("Load test completed - Duration: %.2fs, Success: %d/%d (%.2f%%), Errors: %d, RPS: %.2f",
			testDuration.Seconds(), success, totalRequests, successRate, errors, reqPerSec),
		VirtualUsers:    virtualUsers,
		RequestsPerUser: requestsPerUser,
		TotalRequests:   totalRequests,
		TargetURL:       url,
	})

	// Print summary
	fmt.Println("\n=== Test Summary ===")
	fmt.Printf("Duration: %.2f seconds\n", testDuration.Seconds())
	fmt.Printf("Total Requests: %d\n", totalRequests)
	fmt.Printf("Successful: %d (%.2f%%)\n", success, successRate)
	fmt.Printf("Errors: %d\n", errors)
	fmt.Printf("Requests/sec: %.2f\n", reqPerSec)
	fmt.Printf("\nLogs saved to: %s\n", filename)
}
