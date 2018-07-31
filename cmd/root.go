package main

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
)

// rootCmd is the main "ghpr" command
var rootCmd = &cobra.Command{
	Use:   "ghpr [OPTIONS]",
	Short: "GitHub create a PR tool for command line",
	Run: func(cmd *cobra.Command, args []string) {
		gitCmd := exec.Command("git", "remote", "-v")
		var gitOut bytes.Buffer
		gitCmd.Stdout = &gitOut
		if err := gitCmd.Run(); err != nil {
			log.Fatal(err)
		}
		var repo string
		gitOutLines := strings.Split(gitOut.String(), "\n")
		for _,line := range gitOutLines{
			if strings.Contains(line, Remote){
				afterColon := strings.Split(line, ":")[1]
				repo = strings.Split(afterColon, ".")[0]
				break
			}
		}
		urlStr := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", strings.Split(repo, "/")[0], strings.Split(repo, "/")[1])
		data := url.Values{}
		data.Set("title", args[0])
		var userName string
		if Remote == "upstream"{
			for _,line := range gitOutLines{
				if strings.Contains(line, "origin"){
					// repo = strings.TrimLeft(strings.TrimRight(line, "."), ":")
					afterColon := strings.Split(line, ":")[1]
					userName = strings.Split(afterColon, "/")[0]
					break
				}
			}
			head := fmt.Sprintf("%s:%s",userName,Branch)
			data.Set("head", head)
		} else {
			data.Set("head", Branch)
		}
		data.Set("base", "master")
		client := &http.Client{}
		r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
		// Replace xxxxx with the github personal access token
		r.Header.Add("Authorization", "token xxxxx")
		r.Header.Add("Content-Type", "application/json")

		resp, _ := client.Do(r)
		fmt.Println("Creating a PR.....")
		if resp.Status == "200 OK" {
			fmt.Println("PR created!! :)")
		} else {
			fmt.Println("Ooops, something went wrong :(")
		}
	},
}

var Remote string
var Branch string

func init() {
	f := rootCmd.PersistentFlags()
	f.StringVarP(&Remote, "remote", "r", "upstream", "Remote GitHub repo to make the PR to")
	f.StringVarP(&Branch, "branch", "b", "master", "The branch from which to make the PR from")
}

// Execute executes the command and returns the error
func Execute() error {
	return rootCmd.Execute()
}

