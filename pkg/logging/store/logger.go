package store

import (
	"sync"

	"github.com/beatlabs/patron/log"
)

// Logger is a testing helper implementing the logger interface provided by patron. This struct is used
// in order to capture the logs inside test functions and make assertions on that
type Logger struct {
	*sync.Mutex
	errs [][]interface{}
}

var logger *Logger

// NewLogger returns a new logger
func NewLogger() *Logger {
	lg := &Logger{
		errs:  [][]interface{}{},
		Mutex: &sync.Mutex{},
	}

	lg.Lock()
	logger = lg
	lg.Unlock()
	return lg
}

// GetLogger returns the current logger
func GetLogger() *Logger {
	if logger != nil {
		return logger
	}

	return NewLogger()
}

// FactoryLogger is a function required by patron in order to set the logger
func FactoryLogger(fls map[string]interface{}) log.Logger {
	logger := GetLogger()
	return logger
}

// Sub is a method requred by patron to create a sub logger
func (nl *Logger) Sub(map[string]interface{}) log.Logger {
	return nl
}

// Panic ...
func (nl *Logger) Panic(args ...interface{}) {
	nl.Lock()
	nl.errs = append(nl.errs, append([]interface{}{"panic"}, args...))
	nl.Unlock()
}

// Panicf ...
func (nl *Logger) Panicf(msg string, args ...interface{}) {
	nl.Lock()
	nl.errs = append(nl.errs, append([]interface{}{"panic", msg}, args...))
	nl.Unlock()
}

// Fatal ...
func (nl *Logger) Fatal(args ...interface{}) {
	nl.Lock()
	nl.errs = append(nl.errs, append([]interface{}{"fatal"}, args...))
	nl.Unlock()
}

// Fatalf ...
func (nl *Logger) Fatalf(msg string, args ...interface{}) {
	nl.Lock()
	nl.errs = append(nl.errs, append([]interface{}{"fatal", msg}, args...))
	nl.Unlock()
}

// Error ...
func (nl *Logger) Error(args ...interface{}) {
	nl.Lock()
	nl.errs = append(nl.errs, append([]interface{}{"error"}, args...))
	nl.Unlock()
}

// Errorf ...
func (nl *Logger) Errorf(msg string, args ...interface{}) {
	nl.Lock()
	nl.errs = append(nl.errs, append([]interface{}{"error", msg}, args...))
	nl.Unlock()
}

// Warn ...
func (nl *Logger) Warn(args ...interface{}) {
	nl.Lock()
	nl.errs = append(nl.errs, append([]interface{}{"warn"}, args...))
	nl.Unlock()
}

// Warnf ...
func (nl *Logger) Warnf(msg string, args ...interface{}) {
	nl.Lock()
	nl.errs = append(nl.errs, append([]interface{}{"warn", msg}, args...))
	nl.Unlock()
}

// Info ...
func (nl *Logger) Info(args ...interface{}) {
	nl.Lock()
	nl.errs = append(nl.errs, append([]interface{}{"info"}, args...))
	nl.Unlock()
}

// Infof ...
func (nl *Logger) Infof(msg string, args ...interface{}) {
	nl.Lock()
	nl.errs = append(nl.errs, append([]interface{}{"info", msg}, args...))
	nl.Unlock()
}

// Debug ...
func (nl *Logger) Debug(args ...interface{}) {
	nl.Lock()
	nl.errs = append(nl.errs, append([]interface{}{"debug"}, args...))
	nl.Unlock()
}

// Debugf ...
func (nl *Logger) Debugf(msg string, args ...interface{}) {
	nl.Lock()
	nl.errs = append(nl.errs, append([]interface{}{"debug", msg}, args...))
	nl.Unlock()
}

// GetErrors ...
func (nl *Logger) GetErrors() [][]interface{} {
	nl.Lock()
	res := nl.errs
	nl.Unlock()
	return res
}
