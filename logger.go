package main

import (
	"flag"
	"log"
	"os"
	"sync"
)

type debugLogging bool
type errorLogging bool

var (
	debugLogger debugLogging
	errorLogger errorLogging

	_errorLogger     *log.Logger
	_debugLogger     *log.Logger
	debugLogFilename = "debug.log"

	once sync.Once
)

func init() {
	errorLogger = true
	flag.BoolVar((*bool)(&debugLogger), "debug", false, "log in debug level")

	errorLogFile := os.Stderr
	_errorLogger = log.New(errorLogFile, "[ERROR] ", log.LstdFlags)
}

func (d debugLogging) getDebugLogger() *log.Logger {
	once.Do(func() {
		debugLogFile, err := os.Create(debugLogFilename)
		if err != nil {
			// just panic instead of hiding the problem
			panic(err)
		}
		_debugLogger = log.New(debugLogFile, "[DEBUG] ", log.LstdFlags)
	})
	return _debugLogger
}

func (d debugLogging) Printf(format string, args ...interface{}) {
	if d {
		d.getDebugLogger().Printf(format, args...)
	}
}

func (d debugLogging) Print(args ...interface{}) {
	if d {
		d.getDebugLogger().Print(args...)
	}
}

func (d debugLogging) Println(args ...interface{}) {
	if d {
		d.getDebugLogger().Println(args...)
	}
}

func (e errorLogging) Printf(format string, args ...interface{}) {
	if e {
		_errorLogger.Printf(format, args...)
	}
}

func (e errorLogging) Print(args ...interface{}) {
	if e {
		_errorLogger.Print(args...)
	}
}

func (e errorLogging) Println(args ...interface{}) {
	if e {
		_errorLogger.Println(args...)
	}
}

func (e errorLogging) Panicln(args ...interface{}) {
	if e {
		_errorLogger.Panicln(args...)
	}
}

func (e errorLogging) Fatal(args ...interface{}) {
	if e {
		_errorLogger.Fatal(args...)
	}
}
