package main

import (
	"log"
	"os"
	"time"
)

func setupLogs() *os.File {
	file, err := os.Create("./shared/pharmacy-content/logs/" + time.Now().Format("2006-01-02 15_04_05") + ".log")
	if err != nil {
		log.Fatalf("Error openning log file: %v", err)
	}
	log.SetOutput(file)
	return file
}
