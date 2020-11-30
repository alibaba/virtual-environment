package logger

import (
	"fmt"
	"os"
	"time"
)

var logLevel = LEVEL_INFO

func SetLevel(level string) bool {
	switch level {
	case ERROR:
		logLevel = LEVEL_ERROR
		return true
	case INFO:
		logLevel = LEVEL_INFO
		return true
	case DEBUG:
		logLevel = LEVEL_DEBUG
		return true
	}
	return false
}

func Fatal(msg ...interface{}) {
	fmt.Fprintf(os.Stderr, "[Fatal] %s %s", getTime(), fmt.Sprintln(msg...))
	os.Exit(1)
}

func Error(err error, msg ...interface{}) {
	if err == nil {
		fmt.Fprintf(os.Stderr, "[Error] %s %s", getTime(), fmt.Sprintln(msg...))
	} else {
		fmt.Fprintf(os.Stderr, "[Error] %s {%s} %s", getTime(), err.Error(), fmt.Sprintln(msg...))
	}
}

func Warn(msg ...interface{}) {
	fmt.Fprintf(os.Stderr, "[Warn] %s %s", getTime(), fmt.Sprintln(msg...))
}

func Info(msg ...interface{}) {
	if logLevel >= LEVEL_INFO {
		fmt.Printf("[Info] %s %s", getTime(), fmt.Sprintln(msg...))
	}
}

func Debug(msg ...interface{}) {
	if logLevel >= LEVEL_DEBUG {
		fmt.Printf("[Debug] %s %s", getTime(), fmt.Sprintln(msg...))
	}
}

func getTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
