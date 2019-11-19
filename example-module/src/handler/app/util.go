package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}

// parse the error object and return the exit status, err message and a boolean indicating if able to parse exitStatus
func parseExitError(err error) (exitStatus int, stderr []byte, ok bool) {
	if exitErr, ok := err.(*exec.ExitError); ok {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus(), exitErr.Stderr, true
		}
		return bool2int(!exitErr.Success()), exitErr.Stderr, true
	}
	return 1, nil, false

}

// if there are errors write to stderr and exit
// if no errors exit with exit code set to 0
func writeStderrAndExit(err error) {
	if err == nil {
		os.Exit(0)
	}

	exitStatus, stderr, ok := parseExitError(err)
	if !ok {
		_, err = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if stderr != nil {
		_, err = os.Stderr.Write(stderr)
	}
	os.Exit(exitStatus)
}
