package main

import (
	"log"
	"os"
)

var logLevel = LEVEL_INFO

func initLogger() {
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
