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

func execCmd(cmd string) (returnCode int, output string) {
	debugLogger.Println(cmd)
	argv := strings.Split(cmd, " ")
	bytes, err := exec.Command(argv[0], argv[1:]...).CombinedOutput()
	output = string(bytes)
	if err != nil {
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

	cmd := "git log --graph --color --decorate --all --oneline"
	if len(os.Args) > 1 {
		cmd += " " + strings.Join(os.Args[1:], " ")
	}
	returnCode, output := execCmd(cmd)
	if returnCode != 0 {
		os.Exit(returnCode)
	}
	DrawUI(strings.Split(output, "\n"))
}
