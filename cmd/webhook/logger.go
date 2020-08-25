package main

import (
	"log"
	"os"
	"time"
)

var logLevel = LEVEL_INFO

func initLogger() {
	ticker := time.NewTicker(2 * time.Minute)
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
			logDebug("Ticker:", time.Now().Format("2006-01-02 15:04:05"), ", LogLevel:", logLevel)
		}
	}(ticker)
}

func logError(msg ...interface{}) {
	if logLevel >= LEVEL_ERROR {
		log.Println(msg)
	}
}

func logInfo(msg ...interface{}) {
	if logLevel >= LEVEL_INFO {
		log.Println(msg)
	}
}

func logDebug(msg ...interface{}) {
	if logLevel >= LEVEL_DEBUG {
		log.Println(msg)
	}
}

func logFatal(msg ...interface{}) {
	log.Fatal(msg)
}
