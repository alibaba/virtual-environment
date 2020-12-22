package main

import (
	"alibaba.com/virtual-env-operator/pkg/shared/logger"
	"fmt"
	"github.com/go-logr/logr"
)

type ktLogger struct {
	level      int
	name       string
	parameters []interface{}
}

func (l *ktLogger) Enabled() bool {
	return true
}

func (l *ktLogger) Info(msg string, keysAndValues ...interface{}) {
	if l.level >= logger.LEVEL_INFO {
		logger.Info(l.enrichMessage(msg, keysAndValues))
	}
}

func (l *ktLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	if l.level >= logger.LEVEL_ERROR {
		logger.Error(err, l.enrichMessage(msg, keysAndValues))
	}
}

func (l *ktLogger) WithValues(keysAndValues ...interface{}) logr.Logger {
	return &ktLogger{logger.LEVEL_INFO, "", keysAndValues}
}

func (l *ktLogger) WithName(name string) logr.Logger {
	return &ktLogger{logger.LEVEL_INFO, name, make([]interface{}, 0)}
}

func (l *ktLogger) V(level int) logr.InfoLogger {
	return &ktLogger{level, "", make([]interface{}, 0)}
}

func (l *ktLogger) enrichMessage(msg string, parameters []interface{}) string {
	if l.name != "" {
		if len(l.parameters) > 0 || len(parameters) > 0 {
			return fmt.Sprintf("{ %s:%s%s } %s", l.name, strJoin(l.parameters), strJoin(parameters), msg)
		} else {
			return fmt.Sprintf("{ %s } %s", l.name, msg)
		}
	} else if len(l.parameters) > 0 || len(parameters) > 0 {
		return fmt.Sprintf("{%s%s } %s", strJoin(l.parameters), strJoin(parameters), msg)
	}
	return msg
}

func strJoin(a ...interface{}) string {
	str := ""
	for index := 0; index < len(a); index++ {
		if index > 0 {
			str = str + fmt.Sprintf(", %v", a[index])
		} else {
			str = fmt.Sprintf(" %v", a[index])
		}
	}
	return str
}
