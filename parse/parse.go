package parse

import (
	"fmt"
	"strings"

	"github.com/aswinmprabhu/github-pr-cli/git"
)

func parse(line string) string {
	// check if remote is ssh based
	if strings.Contains(line, "git@github.com") {
		stringAfterColon := strings.Split(line, ":")[1]
		return strings.Split(stringAfterColon, ".")[0]
	}
	stringAfterGithubDomain := strings.Split(line, "https://github.com/")[1]
	return strings.Split(stringAfterGithubDomain, " ")[0]
}

// Remote parses out a remote as username/reponame from "git remote -v"
func Remote(remoteName string) (string, error) {
	gitOut, err := git.GetGitOutput("remote", "-v")
	if err != nil {
		return "", fmt.Errorf("Couldn't parse the remote : %v", err)
	}
	var repo string
	gitOutLines := strings.Split(gitOut, "\n")
	f := 0
	// parse the repo as username/reponame
	for _, line := range gitOutLines {
		// sample gitOutLines elements :
		// ssh : git@github.com:aswinmprabhu/github-pr-cli.git
		// https : https://github.com/aswinmprabhu/github-pr-cli
		if strings.Contains(line, remoteName) {
			repo = parse(line)
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
	gitOut, err := git.GetGitOutput("symbolic-ref", "--short", "HEAD")
	if err != nil {
		return "", fmt.Errorf("Couldn't parse the current branch : %v", err)
	}
	gitOutLines := strings.Split(gitOut, "\n")
	return gitOutLines[0], nil
}
