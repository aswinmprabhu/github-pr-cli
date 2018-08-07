# github-pr-cli
A simple command line utility build using Golang to create GitHub PRs from the comfort of your terminal

```
$ghpr -h
GitHub create a PR tool for command line

Usage:
  ghpr <title> [flags]

Flags:
  -b, --branch string   The branch from which the PR is to be made (default "master")
  -h, --help            help for ghpr
  -r, --remote string   Remote GitHub repo to which the PR is to be made (default "upstream")
```


# Usage

1. Create a config file in the home directory
```
$touch .ghpr.json
```
2. Add the config options and the [access token](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/) to the file
```
{
  "token": "Your token",
  "debug": false
}
```
