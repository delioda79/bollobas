package store

import (
	"github.com/beatlabs/patron/log"
	"sync"
)

type StoreLogger struct {
	*sync.Mutex
	errs [][]interface{}
}

var logger *StoreLogger

func NewLogger() *StoreLogger {
	lg := &StoreLogger{
		errs: [][]interface{}{},
		Mutex: &sync.Mutex{},
	}

	logger = lg
	return lg
}

func GetLogger() *StoreLogger {
	if logger != nil {
		return logger
	}

	return NewLogger()
}

func FactoryLogger(fls map[string]interface{}) log.Logger {
	logger := GetLogger()
	return logger
}

func (nl *StoreLogger) Sub(map[string]interface{}) log.Logger {
	return nl
}

func (nl *StoreLogger) Panic(args ...interface{}) {
	nl.Lock()
	nl.errs = append(nl.errs, append([]interface{}{"panic"}, args...))
	nl.Unlock()
}

func (nl *StoreLogger) Panicf(msg string, args ...interface{}) {
	nl.Lock()
	nl.errs = append(nl.errs, append([]interface{}{"panic", msg}, args...))
	nl.Unlock()
}

func (nl *StoreLogger) Fatal(args ...interface{}) {
	nl.Lock()
	nl.errs = append(nl.errs, append([]interface{}{"fatal"}, args...))
	nl.Unlock()
}

func (nl *StoreLogger) Fatalf(msg string, args ...interface{}) {
	nl.Lock()
	nl.errs = append(nl.errs, append([]interface{}{"fatal", msg}, args...))
	nl.Unlock()
}

func (nl *StoreLogger) Error(args ...interface{}) {
	nl.Lock()
	nl.errs = append(nl.errs, append([]interface{}{"error"}, args...))
	nl.Unlock()
}

func (nl *StoreLogger) Errorf(msg string, args ...interface{}) {
	nl.Lock()
	nl.errs = append(nl.errs, append([]interface{}{"error", msg}, args...))
	nl.Unlock()
}

func (nl *StoreLogger) Warn(args ...interface{}) {
	nl.Lock()
	nl.errs = append(nl.errs, append([]interface{}{"warn"}, args...))
	nl.Unlock()
}

func (nl *StoreLogger) Warnf(msg string, args ...interface{}) {
	nl.Lock()
	nl.errs = append(nl.errs, append([]interface{}{"warn", msg}, args...))
	nl.Unlock()
}

func (nl *StoreLogger) Info(args ...interface{}) {
	nl.Lock()
	nl.errs = append(nl.errs, append([]interface{}{"info"}, args...))
	nl.Unlock()
}

func (nl *StoreLogger) Infof(msg string, args ...interface{}) {
	nl.Lock()
	nl.errs = append(nl.errs, append([]interface{}{"info", msg}, args...))
	nl.Unlock()
}

func (nl *StoreLogger) Debug(args ...interface{}) {
	nl.Lock()
	nl.errs = append(nl.errs, append([]interface{}{"debug"}, args...))
	nl.Unlock()
}

func (nl *StoreLogger) Debugf(msg string, args ...interface{}) {
	nl.Lock()
	nl.errs = append(nl.errs, append([]interface{}{"debug", msg}, args...))
	nl.Unlock()
}

func (nl *StoreLogger) GetErrors() [][]interface{} {
	return nl.errs
}
