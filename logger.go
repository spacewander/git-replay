package main

import (
	"flag"
	"io"
	"log"
	"os"
)

type debugLogging bool
type errorLogging bool

var (
	debugLogger debugLogging
	errorLogger errorLogging

	logFile io.Writer

	_errorLogger = log.New(os.Stdout, "[ERROR] ", log.LstdFlags)
	_debugLogger = log.New(os.Stdout, "[DEBUG] ", log.LstdFlags)
)

func init() {
	errorLogger = true
	flag.BoolVar((*bool)(&debugLogger), "debug", false, "log in debug level")
}

func initLog() {
	logFile = os.Stderr
	log.SetOutput(logFile)
	_errorLogger = log.New(logFile, "[ERROR] ", log.LstdFlags)
	_debugLogger = log.New(logFile, "[DEBUG] ", log.LstdFlags)
}

func (d debugLogging) Printf(format string, args ...interface{}) {
	if d {
		_debugLogger.Printf(format, args...)
	}
}

func (d debugLogging) Print(args ...interface{}) {
	if d {
		_debugLogger.Print(args...)
	}
}

func (d debugLogging) Println(args ...interface{}) {
	if d {
		_debugLogger.Println(args...)
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
