package main

import (
	"log"
	"os"
	"time"
)

var logLevel = LEVEL_INFO

func initLogger() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	go func(t *time.Ticker) {
		for {
			<-t.C
			switch os.Getenv(CONF_LOG_LEVEL) {
			case ERROR:
				logLevel = LEVEL_ERROR
				break
			case INFO:
				logLevel = LEVEL_INFO
				break
			case DEBUG:
				logLevel = LEVEL_DEBUG
				break
			}
			log.Println(logLevel, " - Ticker: ", time.Now().Format("2006-01-02 15:04:05"))
		}
	}(ticker)
}
