package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var (
	printVersion bool
	VERSION      string
)

func init() {
	flag.BoolVar((*bool)(&printVersion), "version", false, "print version")
}

func sh(cmd string) (returnCode int, output string) {
	debugLogger.Println(cmd)
	bytes, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	output = string(bytes)
	if err != nil || strings.HasSuffix(output, "command not found\n") {
		errorLogger.Print(output)
		exitErr, ok := err.(*exec.ExitError)
		if ok {
			errorLogger.Print(exitErr)
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				returnCode = status.ExitStatus()
				return returnCode, output
			}
		}
		return 1, output
	}
	return 0, strings.TrimSpace(output)
}

func main() {
	flag.Parse()

	if printVersion {
		// use `go build -ldflags "-X main.VERSION=$(git rev-parse HEAD)"`
		// to inject the commit sha1 as version
		fmt.Println("Version: ", VERSION)
		os.Exit(0)
	}

	returnCode, output := sh("git log --graph --color --decorate --all --date=iso --oneline")
	if returnCode != 0 {
		os.Exit(returnCode)
	}
	DrawUI(output)
}
