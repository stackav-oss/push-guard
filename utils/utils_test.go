// Copyright 2025 Stack AV Co.
// SPDX-License-Identifier: Apache-2.0
package utils

import (
	"errors"
	"runtime"
	"strconv"
	"testing"
)

type MockExecutor struct {
	ExecutableFunc func() (string, error)
	LookPathFunc   func(string) (string, error)
}

func (e *MockExecutor) Executable() (string, error) {
	return e.ExecutableFunc()
}

func (e *MockExecutor) LookPath(file string) (string, error) {
	return e.LookPathFunc(file)
}

// TestIsSafeRemoteParam tests the IsSafeRemote function to ensure it correctly identifies
// whether a given remote URL is considered safe or not. It runs a series of test cases
// with different remote URLs and their expected safety status. Additionally, it tests
// the opposite cases to ensure the function's accuracy in both positive and negative scenarios.
func TestIsSafeRemoteParam(t *testing.T) {
	ProtocolAndDomainAllowList = []string{
		"github.com:",
		"git@github.com:",
		"https://github.com/",
		"ssh://git@github.com/",
	}
	DirectoryAllowList = []string{
		"allowed-organization",
		"allowed-user-account",
	}
	DirectoryRegexAllowList = []string{
		"^[a-z,0-9]+_allowed_suffix/",
	}
	tests := []struct {
		remote string
		want   bool
	}{
		{"github.com:blocked-organization/test.git", false},                // allowed: domain, blocked: organization
		{"git@github.com:blocked-organization/test.git", false},            // allowed: domain, blocked: organization
		{"gitlab.com:allowed-organization/test.git", false},                // allowed: organization, blocked: domain
		{"git@gitlab.com:allowed-organization/test.git", false},            // allowed: organization, blocked: protocol + domain
		{"http://gitlab.com/allowed-organization/test.git", false},         // allowed: organization, blocked: protocol + domain
		{"https://gitlab.com/allowed-organization/test.git", false},        // allowed: organization, blocked: domain
		{"git@bitbucket.org:allowed-organization/test.git", false},         // allowed: organization, blocked: domain
		{"https://bitbucket.org/allowed-organization/test.git", false},     // allowed: organization, blocked: domain
		{"https://github.com/allowed-user-account/test.git", true},         // allowed: protocol + domain + user account + repo.git
		{"https://github.com/allowed-organization/test", true},             // allowed: protocol + domain + organization + repo without .git
		{"git@github.com:blocked-organization/test.git", false},            // allowed: protocol + domain, blocked: organization
		{"github.com:allowed-organization/test.git", true},                 // allowed: domain + organization
		{"github.com:allowed-user-account/test.git", true},                 // allowed: domain + user account
		{"http://github.com/account123_allowed_suffix/test.git", false},    // allowed: domain + user account suffix, blocked: protocol
		{"https://github.com/account123_allowed_suffix/test.git", true},    // allowed: protocol + domain + user account suffix
		{"git@github.com:account123_allowed_suffix/test.git", true},        // allowed: protocol + domain + user account suffix
		{"https://github.com/account123_allowed_suffix_test.git", false},   // allowed: protocol + domain, invalid: repo path
		{"git@github.com:account123/test.git", false},                      // allowed: protocol + domain, blocked: user account suffix
		{"git@github.com:_allowed_suffix/test.git", false},                 // allowed: protocol + domain, blocked: user account suffix
		{"git@github.com:allowed-organization/*.git", true},                // allowed: protocol + domain + organization + wild card
		{"ssh://git@github.com/allowed-organization/test.git", true},       // allowed: protocol + domain + organization
		{"https://github.com/allowed-organization/private-repo.git", true}, // allowed: protocol + domain + organization, different repo name
		{"git@github.com:allowed-organization/private-repo.git", true},     // allowed: protocol + domain + organization, different repo name
	}

	for _, tt := range tests {
		t.Run(tt.remote, func(t *testing.T) {
			got := IsSafeRemote(tt.remote)
			if got != tt.want {
				t.Errorf("got %q, want %q", strconv.FormatBool(got), strconv.FormatBool(tt.want))
			}
		})
	}

	// Test the opposite cases
	for _, tt := range tests {
		t.Run("Not "+tt.remote, func(t *testing.T) {
			got := !IsSafeRemote(tt.remote)
			if got != !tt.want {
				t.Errorf("got %q, want %q", strconv.FormatBool(got), strconv.FormatBool(!tt.want))
			}
		})
	}
}

func TestIsSafeRemoteEdgeCases(t *testing.T) {
	ProtocolAndDomainAllowList = []string{
		"github.com:",
		"git@github.com:",
		"https://github.com/",
		"ssh://git@github.com/",
	}
	DirectoryAllowList = []string{
		"allowed-organization",
		"allowed-user-account",
	}
	cases := []struct {
		remote string
		want   bool
	}{
		{"", false},
		{"ftp://github.com/allowed-organization/test.git", false},
		{"github.com", false}, // no colon or scheme
	}
	for _, c := range cases {
		got := IsSafeRemote(c.remote)
		if got != c.want {
			t.Errorf("IsSafeRemote(%q) = %v, want %v", c.remote, got, c.want)
		}
	}
}

func TestRemovePath(t *testing.T) {
	path := "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
	got := RemovePath(path, "/usr/local/bin")
	want := "/usr/local/sbin:/usr/sbin:/usr/bin:/sbin:/bin"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	got = RemovePath(path, "")
	want = path
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRemovePathEdgeCases(t *testing.T) {
	longPath := "/usr/local/bin:/usr/local/sbin:/usr/bin:/bin:/usr/sbin:/sbin"
	cases := []struct {
		name     string
		path     string
		remove   string
		expected string
	}{
		{"RemoveFirst", longPath, "/usr/local/bin", "/usr/local/sbin:/usr/bin:/bin:/usr/sbin:/sbin"},
		{"RemoveMiddle", longPath, "/usr/bin", "/usr/local/bin:/usr/local/sbin:/bin:/usr/sbin:/sbin"},
		{"NotPresent", longPath, "/fake/path", longPath},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := RemovePath(c.path, c.remove)
			if got != c.expected {
				t.Errorf("RemovePath(%q, %q) = %q, want %q", c.path, c.remove, got, c.expected)
			}
		})
	}
}

func TestLocateGitBinary(t *testing.T) {
	ExecutableFuncGood := func() (string, error) {
		return "/current/running/git", nil
	}
	ExecutableFuncBad := func() (string, error) {
		return "", errors.New("Not found")
	}
	LookPathFuncGood := func(file string) (string, error) {
		return "/default/installed/git", nil
	}
	LookPathFuncBad := func(file string) (string, error) {
		return "", errors.New("Not found")
	}
	mockExec := &MockExecutor{}
	mockExec.ExecutableFunc = ExecutableFuncGood
	mockExec.LookPathFunc = LookPathFuncGood
	got := LocateGitBinary(mockExec)
	want := "/default/installed/git"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	mockExec.ExecutableFunc = ExecutableFuncBad
	got = LocateGitBinary(mockExec)
	want = DefaultGitCommand(runtime.GOOS)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	mockExec.ExecutableFunc = ExecutableFuncGood
	mockExec.LookPathFunc = LookPathFuncBad
	got = LocateGitBinary(mockExec)
	want = DefaultGitCommand(runtime.GOOS)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestLocateGitBinaryEdgeCases(t *testing.T) {
	mockExec := &MockExecutor{
		ExecutableFunc: func() (string, error) { return "/custom/git", nil },
		LookPathFunc:   func(file string) (string, error) { return "/usr/bin/git", nil },
	}
	// Confirm that both Executable() and LookPath() succeed
	got := LocateGitBinary(mockExec)
	if got != "/usr/bin/git" {
		t.Errorf("LocateGitBinary() = %q, want %q", got, "/usr/bin/git")
	}

	// Confirm fallback in case both fail
	mockExec.ExecutableFunc = func() (string, error) { return "", errors.New("no exec found") }
	mockExec.LookPathFunc = func(file string) (string, error) { return "", errors.New("no path found") }
	got = LocateGitBinary(mockExec)
	want := DefaultGitCommand(runtime.GOOS)
	if got != want {
		t.Errorf("LocateGitBinary() fallback = %q, want %q", got, want)
	}
}

func TestFindUnsafeRemoteFromStdErr(t *testing.T) {
	ProtectedBranches = []string{"protected"} // mock protected branches
	ProtocolAndDomainAllowList = []string{"github.com:"}
	DirectoryAllowList = []string{"allowed-organization"}
	lines := []string{
		"error: src refspec push does not match any",
		"error: failed to push some refs to 'git'",
	}
	got_result, got_remote := FindUnsafeRemoteFromStdErr(lines, 1)
	want_result, want_remote := false, ""
	if got_result != want_result {
		t.Errorf("got %q, want %q", strconv.FormatBool(got_result), strconv.FormatBool(want_result))
	}
	if got_remote != want_remote {
		t.Errorf("got %q, want %q", got_remote, want_remote)
	}
	lines = []string{
		"To github.com:test.git",
		"   f857dfa7..23092a82  test-branch -> test-branch",
	}
	got_result, got_remote = FindUnsafeRemoteFromStdErr(lines, 0)
	want_result, want_remote = true, "github.com:test.git"
	if got_result != want_result {
		t.Errorf("got %q, want %q", strconv.FormatBool(got_result), strconv.FormatBool(want_result))
	}
	if got_remote != want_remote {
		t.Errorf("got %q, want %q", got_remote, want_remote)
	}
	lines = []string{
		"To github.com:allowed-organization/eie-swe.git",
		"   f857dfa7..ec8abedf  protected -> protected",
	}
	got_result, got_remote = FindUnsafeRemoteFromStdErr(lines, 0)
	want_result, want_remote = true, "protected"
	if got_result != want_result {
		t.Errorf("got %q, want %q", strconv.FormatBool(got_result), strconv.FormatBool(want_result))
	}
	if got_remote != want_remote {
		t.Errorf("got %q, want %q", got_remote, want_remote)
	}
}
