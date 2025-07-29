//go:build darwin || linux

package platform

import (
	"os"
	"push-guard/utils"
	"syscall"
)

func ExecuteGit(gitBinaryPath string) {
	utils.Logger.Debug("ExecuteGit", "gitBinaryPath", gitBinaryPath, "Args", os.Args)
	execErr := syscall.Exec(gitBinaryPath, os.Args, os.Environ())
	if execErr != nil {
		panic(execErr)
	}
}
