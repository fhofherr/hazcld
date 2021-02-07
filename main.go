package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/fhofherr/hazcld/internal/process"
)

// Possible exit codes.
const (
	ExitFound = iota
	ExitNotFound
	ExitErr
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage:\n\n%s <regexp> <pid>\n", os.Args[0])
		os.Exit(ExitErr)
	}

	re, err := regexp.Compile(os.Args[1])
	if err != nil {
		fmt.Printf("Compile regexp: %v\n", err)
		os.Exit(ExitErr)
	}
	pid, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Get pid: %v\n", err)
		os.Exit(ExitErr)
	}

	found, err := process.HasChildProcess(pid, re)
	if err != nil {
		fmt.Printf("Search child processes: %v\n", err)
		os.Exit(ExitErr)
	}
	if !found {
		os.Exit(ExitNotFound)
	}
	os.Exit(ExitFound)
}
