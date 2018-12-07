// Package git consists of all the utility functions
// that need to execute the git command.
package git

import (
	"bytes"
	"fmt"
	"os/exec"
)

// GetGitOutput returns the output of "git remote -v"
func GetGitOutput(args ...string) (string, error) {
	// exec "git remote -v" to get the remotes
	gitCmd := exec.Command("git", args...)
	var gitOut bytes.Buffer
	gitCmd.Stdout = &gitOut
	if err := gitCmd.Run(); err != nil {
		return "", fmt.Errorf("Couldn't get the remotes : %v", err)
	}
	return gitOut.String(), nil
}
