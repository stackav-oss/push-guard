// Copyright 2025 Stack AV Co.
// SPDX-License-Identifier: Apache-2.0
//go:build windows

package platform

import (
	"os"
	"os/exec"
	"push-guard/utils"
	"syscall"

	"golang.org/x/sys/windows"
)

func ExecuteGit(gitBinaryPath string) {
	var cmdArgs []string

	if len(os.Args) > 1 {
		cmdArgs = os.Args[1:]
	}

	cmd := exec.Command(gitBinaryPath, cmdArgs...)
	utils.Logger.Debug("ExecuteGit", "gitBinaryPath", gitBinaryPath, "Args", cmdArgs)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: windows.CREATE_NEW_PROCESS_GROUP,
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	cmd.Wait()
}
