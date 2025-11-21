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

// Global log file handle
var logFile *os.File

func InitLogger() error {
	var err error
	logFile, err = os.Create("requests.log")
	return err
}

func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}

func WriteLog(entry LogEntry) {
	if logFile == nil {
		return
	}

	jsonBytes, _ := json.Marshal(entry)
	logFile.Write(jsonBytes)
	logFile.Write([]byte("\n"))
}
