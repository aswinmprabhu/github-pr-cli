// Package git consists of all the utility functions
// that need to execute the git command.
package git

import (
	"bytes"
	"fmt"
	"os/exec"
)

// GetRemotes returns the output of "git remote -v"
func GetRemotes() (string, error) {
	// exec "git remote -v" to get the remotes
	gitCmd := exec.Command("git", "remote", "-v")
	var gitOut bytes.Buffer
	gitCmd.Stdout = &gitOut
	if err := gitCmd.Run(); err != nil {
		return "", fmt.Errorf("Couldn't get the remotes : %v", err)
	}
	return gitOut.String(), nil
}

// GetCurrentBranch returns the output of "git symbolic-ref --short HEAD"
func GetCurrentBranch() (string, error) {
	// exec "git remote -v" to get the remotes
	gitCmd := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	var gitOut bytes.Buffer
	gitCmd.Stdout = &gitOut
	if err := gitCmd.Run(); err != nil {
		return "", fmt.Errorf("Couldn't get the branch : %v", err)
	}
	return gitOut.String(), nil

}
