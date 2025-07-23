package main

import (
	"fmt"
	"os"
	"push-guard/config"
	"push-guard/platform"
	"push-guard/utils"
	"slices"
)

func HandleDisclaimer(remote string) {
	fmt.Printf("\n\n")
	fmt.Println(config.DecodeConfigString(config.Disclaimer))
	fmt.Println()
	fmt.Printf("Git Remote: %q\n", remote)
	if !utils.Confirmation("Do you want to continue with the action?") {
		fmt.Println("Exiting...")
		utils.SendMessage(remote, false)
		os.Exit(0)
	} else {
		utils.SendMessage(remote, true)
	}
}

func HandleDisclaimerForProtectedBranch(branch string) {
	warning := fmt.Sprintf("You are attempting to push to the protected branch: \"%s\"", branch)
	fmt.Println(warning)
	if !utils.Confirmation("Do you want to continue with the action?") {
		fmt.Println("Exiting...")
		os.Exit(0)
	}
}

func main() {
	gitBinaryPath := utils.LocateGitBinary(&utils.RealExecutor{})
	if len(os.Args) > 1 {
		if os.Args[1] == "--push-guard-version" {
			fmt.Println(config.DecodeConfigString(config.PushGuardVersion))
			os.Exit(0)
		}
		if unsafe, remote := utils.FindUnsafeRemote(gitBinaryPath, os.Args[1:]); unsafe {
			if slices.Contains(utils.ProtectedBranches, remote) {
				HandleDisclaimerForProtectedBranch(remote)
			} else {
				HandleDisclaimer(remote)
			}
		}
	}
	utils.Logger.Debug("Calling ExecuteGit")
	platform.ExecuteGit(gitBinaryPath)
}
