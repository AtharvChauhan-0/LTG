package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type LogEntry struct {
	Timestamp  string `json:"timestamp"`
	Level      string `json:"level"`
	VU         int    `json:"vu"`
	Status     int    `json:"status"`
	DurationMS int64  `json:"duration_ms"`
	Error      string `json:"error,omitempty"`
	URL        string `json:"url"`
	Method     string `json:"method"`
	TestRun    string `json:"test_run"`
}

type TestMetadata struct {
	Timestamp       string `json:"timestamp"`
	Level           string `json:"level"`
	Message         string `json:"message"`
	TestRun         string `json:"test_run"`
	VirtualUsers    int    `json:"virtual_users,omitempty"`
	RequestsPerUser int    `json:"requests_per_user,omitempty"`
	TotalRequests   int    `json:"total_requests,omitempty"`
	TargetURL       string `json:"target_url,omitempty"`
}

var (
	logChan   chan interface{}
	logFile   *os.File
	testRunID string
	encoder   *json.Encoder
)

// InitLogger creates a new timestamped log file and initializes the logger
func InitLogger() (string, error) {
	var err error

	// Generate unique test run ID with timestamp
	testRunID = fmt.Sprintf("test_%s", time.Now().Format("20060102_150405"))

	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		return "", err
	}

	// Create new log file with timestamp
	filename := fmt.Sprintf("logs/loadtest_%s.log", time.Now().Format("20060102_150405"))
	logFile, err = os.Create(filename)
	if err != nil {
		return "", err
	}

	// Buffered channel (async logging)
	logChan = make(chan interface{}, 10000)

	// Create JSON encoder
	encoder = json.NewEncoder(logFile)

	// Logger goroutine
	go func() {
		for entry := range logChan {
			encoder.Encode(entry)
		}
	}()

	return filename, nil
}

// WriteLog writes a request log entry
func WriteLog(entry LogEntry) {
	if logChan != nil {
		// Add test run ID
		entry.TestRun = testRunID
		// Format timestamp in RFC3339 (ISO 8601) for better compatibility
		entry.Timestamp = time.Now().Format(time.RFC3339)
		logChan <- entry
	}
}

// WriteMetadata writes test metadata (start, end, summary)
func WriteMetadata(metadata TestMetadata) {
	if logChan != nil {
		metadata.TestRun = testRunID
		metadata.Timestamp = time.Now().Format(time.RFC3339)
		logChan <- metadata
	}
}

// CloseLogger gracefully shuts down the logger
func CloseLogger() {
	// Stop accepting new entries
	if logChan != nil {
		close(logChan)
	}

	// Close file AFTER all remaining logs finish writing
	if logFile != nil {
		logFile.Close()
	}
}
