package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// ParseRemote parses out a remote as username/reponame from "git remote -v"
func ParseRemote(remoteName string) (string, error) {
	// exec "git remote -v" to get the remotes
	gitCmd := exec.Command("git", "remote", "-v")
	var gitOut bytes.Buffer
	gitCmd.Stdout = &gitOut
	if err := gitCmd.Run(); err != nil {
		return "", fmt.Errorf("Couldn't parse the remote : %v", err)
	}
	var repo string
	gitOutLines := strings.Split(gitOut.String(), "\n")
	f := 0
	// parse the repo as username/reponame
	for _, line := range gitOutLines {
		// sample gitOutLines element : git@github.com:aswinmprabhu/github-pr-cli.git
		if strings.Contains(line, remoteName) {
			stringAfterColon := strings.Split(line, ":")[1]
			repo = strings.Split(stringAfterColon, ".")[0]
			f = 1
			break
		}
	}
	if f == 0 {
		return "", fmt.Errorf("Couldn't parse the remote : Remote not found")
	}
	return repo, nil
}

// CurrentBranch returns the current branch as a string
func CurrentBranch() (string, error) {
	// exec "git remote -v" to get the remotes
	gitCmd := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	var gitOut bytes.Buffer
	gitCmd.Stdout = &gitOut
	if err := gitCmd.Run(); err != nil {
		return "", fmt.Errorf("Couldn't find the current branch : %v", err)
	}
	gitOutLines := strings.Split(gitOut.String(), "\n")
	return gitOutLines[0], nil
}
