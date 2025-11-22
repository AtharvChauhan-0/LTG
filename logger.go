package main

import (
	"encoding/json"
	"os"
	"time"
)

type LogEntry struct {
	Timestamp  time.Time `json:"timestamp"`
	VU         int       `json:"vu"`
	Status     int       `json:"status"`
	DurationMS int64     `json:"duration_ms"`
	Error      string    `json:"error,omitempty"`
	URL        string    `json:"url"`
	Method     string    `json:"method"`
}

var (
	logChan chan LogEntry
	logFile *os.File
)

func InitLogger() error {
	var err error

	// Create or overwrite existing log file
	logFile, err = os.Create("requests.log")
	if err != nil {
		return err
	}

	// Buffered channel (async logging)
	logChan = make(chan LogEntry, 10000)

	// Logger goroutine
	go func() {
		for entry := range logChan {
			jsonBytes, _ := json.Marshal(entry)
			logFile.Write(jsonBytes)
			logFile.Write([]byte("\n"))
		}
	}()

	return nil
}

func WriteLog(entry LogEntry) {
	if logChan != nil {
		logChan <- entry
	}
}

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
