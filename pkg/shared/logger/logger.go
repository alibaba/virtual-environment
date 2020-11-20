package logger

import (
	"log"
	"time"
)

var logLevel = LEVEL_INFO

func SetLevel(level string) {
	switch level {
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

func Fatal(msg ...interface{}) {
	log.Fatal("[Fatal] ", getTime(), msg)
}

func Error(err error, msg ...interface{}) {
	log.Println("[Error] ", getTime(), msg, ": ", err.Error())
}

func Warn(msg ...interface{}) {
	log.Println("[Warn] ", getTime(), msg)
}

func Info(msg ...interface{}) {
	if logLevel >= LEVEL_INFO {
		log.Println("[Info] ", getTime(), msg)
	}
}

func Debug(msg ...interface{}) {
	if logLevel >= LEVEL_DEBUG {
		log.Println("[Debug] ", getTime(), msg)
	}
}

func getTime() string {
	return time.Now().Format("2006-01-02T15:04:05 ")
}
