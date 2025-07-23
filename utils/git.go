package utils

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"push-guard/config"
	"regexp"
	"runtime"
	"slices"
	"strings"
)

var ProtocolAndDomainAllowList []string = strings.Split(config.ProtocolAndDomainAllowList, ",")
var DirectoryAllowList []string = strings.Split(config.DirectoryAllowList, ",")
var DirectoryRegexAllowList []string = strings.Split(config.DirectoryRegexAllowList, ",")
var ProtectedBranches []string = strings.Split(config.ProtectedBranches, ",")

type Executor interface {
	Executable() (string, error)
	LookPath(file string) (string, error)
}

type RealExecutor struct{}

func (e *RealExecutor) Executable() (string, error) {
	return os.Executable()
}

func (e *RealExecutor) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

func IsSafeRemote(url string) bool {
	slog.Debug("DEBUG MESSAGE from IsSafeRemote")
	for _, protocolAndDomain := range ProtocolAndDomainAllowList {
		suffix, hasPrefix := strings.CutPrefix(url, protocolAndDomain)
		if hasPrefix {
			for _, directory := range DirectoryAllowList {
				if strings.HasPrefix(suffix, directory) {
					return true
				}
			}
			for _, regexString := range DirectoryRegexAllowList {
				match, _ := regexp.MatchString(regexString, suffix)
				if match {
					return true
				}
			}
		}
	}
	return false
}

// Provides best guess default locations where the git binary
// might be installed for each of the three platforms.
// This should be used as a last resort if lookpath fails to find git.
func DefaultGitCommand(os string) string {
	switch os {
	case "darwin":
		return "/usr/bin/git"
	case "linux":
		return "/usr/bin/git"
	case "windows":
		return "c:\\Program Files\\Git\\cmd\\git.exe"
	default:
		return "git"
	}
}

// Remove target from envPath
func RemovePath(envPath, target string) string {
	if envPath != "" {
		var newPath []string
		pathList := strings.Split(envPath, string(os.PathListSeparator))
		for _, p := range pathList {
			if p != target {
				newPath = append(newPath, p)
			}
		}
		return strings.Join(newPath, string(os.PathListSeparator))
	}
	return ""
}

// Find the real git binary on the system and return the path to it.
func LocateGitBinary(execExecutor Executor) string {
	exePath, err := execExecutor.Executable()
	if err != nil {
		return DefaultGitCommand(runtime.GOOS)
	}
	Logger.Debug("LocateGitBinary", "exePath", exePath)
	origPath := os.Getenv("PATH")
	newPath := RemovePath(origPath, filepath.Dir(exePath))
	Logger.Debug("LocateGitBinary", "newPath", newPath)
	os.Unsetenv("PATH")
	os.Setenv("PATH", newPath)
	gitPath, err := execExecutor.LookPath("git") // Try to locate the Git binary in the modified PATH
	if err != nil {
		return DefaultGitCommand(runtime.GOOS)
	}
	os.Unsetenv("PATH")
	os.Setenv("PATH", origPath) // Restore the original path
	Logger.Debug("LocateGitBinary", "gitPath", gitPath)
	return gitPath
}

func FindUnsafeRemoteFromStdErr(lines []string, exitCode int) (bool, string) {
	if exitCode == 0 {
		if strings.HasPrefix(lines[0], "To ") {
			remote := strings.Fields(lines[0])[1]
			if !IsSafeRemote(remote) {
				Logger.Debug("FindUnsafeRemote: Unsafe Remote Found", "remote", remote)
				return true, remote
			}
		}
		for _, line := range lines {
			for _, branch := range ProtectedBranches {
				if strings.HasSuffix(line, fmt.Sprintf(" -> %s", branch)) {
					return true, branch
				}
			}
		}
	} else {
		for _, line := range lines {
			if strings.HasPrefix(line, "error: src refspec push does not match any") {
				return false, ""
			}
			if strings.HasPrefix(line, "error: failed to push some refs to ") {
				remote := strings.Split(line, "'")[1]
				if !IsSafeRemote(remote) {
					Logger.Debug("FindUnsafeRemote: Unsafe Remote Found", "remote", remote)
					return true, remote
				}
			}
		}
	}
	return false, ""
}

func FindUnsafeRemote(gitBinaryPath string, args []string) (bool, string) {
	if slices.Contains(args, "push") {
		dry_run_args := append(args, "--dry-run")
		_, stdErr, exitCode := RunCommand(gitBinaryPath, dry_run_args...)
		if stdErr != "" {
			Logger.Debug("FindUnsafeRemote", "dry run output", stdErr)
			lines := strings.FieldsFunc(stdErr, func(c rune) bool { return c == '\n' || c == '\r' })
			return FindUnsafeRemoteFromStdErr(lines, exitCode)
		}
	}
	return false, ""
}
