# github-pr-cli
A simple command line utility build using Golang to create GitHub PRs from the comfort of your terminal

```shell
$ghpr -h
GitHub create a PR tool for command line

Usage:
  ghpr [OPTIONS] <title> [flags]

Flags:
  -b, --branch string   The branch from which to make the PR from (default "master")
  -h, --help            help for ghpr
  -r, --remote string   Remote GitHub repo to make the PR to (default "upstream")
```

## TODO
1. Add viper config
