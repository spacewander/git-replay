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

	scriptName string
)

func init() {
	flag.BoolVar(&printVersion, "version", false, "print version")
	flag.StringVar(&scriptName, "script", "", "specific lua script to run")
}

func panicIfScriptIsInvalid(filename string) {
	if stat, err := os.Stat(filename); err != nil {
		errorLogger.Panicln(err)
	} else if !stat.Mode().IsRegular() {
		errorLogger.Panicln(filename + " is not a regular file")
	}
}

func execCmd(argv []string) (returnCode int, output string) {
	debugLogger.Println(argv)
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
	if scriptName != "" {
		panicIfScriptIsInvalid(scriptName)
	}

	returnCode, path := execCmd(strings.Split(`git rev-parse --show-toplevel`, " "))
	if returnCode != 0 {
		os.Exit(returnCode)
	}
	if err := InitRepo(path); err != nil {
		errorLogger.Fatal(err)
	}

	cmd := strings.Split(`git log --graph --all --color`, " ")
	cmd = append(cmd, `--format=%C(yellow)%H%Creset%C(auto)%d %s`)
	args := flag.Args()
	cmd = append(cmd, args...)
	returnCode, output := execCmd(cmd)
	if returnCode != 0 {
		os.Exit(returnCode)
	}
	DrawUI(strings.Split(output, "\n"))
}
