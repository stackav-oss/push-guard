// Copyright 2025 Stack AV Co.
// SPDX-License-Identifier: Apache-2.0
package utils

import (
	"bytes"
	"os"
	"os/exec"
	"syscall"
)

func RunCommand(cmd string, args ...string) (stdout string, stderr string, exitCode int) {
	Logger.Debug("RunCommand", "cmd", cmd, "args", args)
	var outbuf, errbuf bytes.Buffer
	command := exec.Command(cmd, args...)
	command.Stdout = &outbuf
	command.Stderr = &errbuf
	err := command.Run()
	stdout = outbuf.String()
	stderr = errbuf.String()
	if err != nil {
		Logger.Debug("RunCommand", "error", err)
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus := exitError.Sys().(syscall.WaitStatus)
			exitCode = waitStatus.ExitStatus()
		} else {
			exitCode = 1
		}
	}
	Logger.Debug("RunCommand Exit")
	return
}

func RunCommandStreamStdoutStderr(cmd string, args ...string) int {
	command := exec.Command(cmd, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus := exitError.Sys().(syscall.WaitStatus)
			return waitStatus.ExitStatus()
		}
		return 1
	}
	return 0
}
