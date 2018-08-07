package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// rootCmd is the main "ghpr" command
var rootCmd = &cobra.Command{
	Use:   "ghpr [OPTIONS] <title>",
	Short: "GitHub create a PR tool for command line",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("PR title required")
			os.Exit(0)
		}
		if Debug {
			fmt.Println("Remote:", Remote, "Branch:", Branch, "Title:", args[0])
		}
		gitCmd := exec.Command("git", "remote", "-v")
		var gitOut bytes.Buffer
		gitCmd.Stdout = &gitOut
		if err := gitCmd.Run(); err != nil {
			fmt.Println("Not a git repo")
			os.Exit(0)
		}
		var repo string
		gitOutLines := strings.Split(gitOut.String(), "\n")
		f := 0
		for _, line := range gitOutLines {
			if strings.Contains(line, Remote) {
				afterColon := strings.Split(line, ":")[1]
				repo = strings.Split(afterColon, ".")[0]
				f = 1
				break
			}
		}
		if f == 0 {
			fmt.Println("Remote not found")
			os.Exit(0)
		}
		urlStr := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", strings.Split(repo, "/")[0], strings.Split(repo, "/")[1])
		jsonValues := map[string]string{"title": args[0], "base": "master"}
		var userName string
		if Remote == "upstream" {
			for _, line := range gitOutLines {
				if strings.Contains(line, "origin") {
					afterColon := strings.Split(line, ":")[1]
					userName = strings.Split(afterColon, "/")[0]
					break
				}
			}
			head := fmt.Sprintf("%s:%s", userName, Branch)
			jsonValues["head"] = head
		} else {
			jsonValues["head"] = Branch
		}
		jsonObj, _ := json.Marshal(jsonValues)
		client := &http.Client{}
		r, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer(jsonObj)) // URL-encoded payload
		AuthVal := fmt.Sprintf("token %s", Token)
		r.Header.Add("Authorization", AuthVal)
		r.Header.Add("Content-Type", "application/json")

		resp, err := client.Do(r)
		if err != nil {
			log.Fatal(err)
		}
		// defer resp.Body.Close()
		fmt.Println("Creating a PR.....")
		if Debug {
			bytes, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(string(bytes))
		}
		if resp.Status == "201 Created" {
			fmt.Println("PR created!! :)")
		} else {
			fmt.Println("Ooops, something went wrong :(")
		}
	},
}

var Remote string
var Branch string
var Token string
var Debug bool

func init() {
	cfgFile := fmt.Sprintf("%s/.ghpr.json", os.Getenv("HOME"))
	viper.SetConfigFile(cfgFile)
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	Debug = viper.GetBool("debug")
	fmt.Println(Token)
	Token = viper.GetString("token")
	if Debug {
		fmt.Println("Debug:", Debug, "Token:", Token)
	}
	f := rootCmd.PersistentFlags()
	f.StringVarP(&Remote, "remote", "r", "upstream", "Remote GitHub repo to which the PR is to be made")
	f.StringVarP(&Branch, "branch", "b", "master", "The branch from which the PR is to be made")
}

// Execute executes the command and returns the error
func Execute() error {
	return rootCmd.Execute()
}
