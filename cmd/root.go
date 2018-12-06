package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/aswinmprabhu/github-pr-cli/browser"
	"github.com/aswinmprabhu/github-pr-cli/editor"
	"github.com/aswinmprabhu/github-pr-cli/parse"
	"github.com/aswinmprabhu/github-pr-cli/request"
)

// rootCmd is the main "ghpr" command
var rootCmd = &cobra.Command{
	Use:   "ghpr <title>",
	Short: "Create github pull requests from the command line",
	Run: func(cmd *cobra.Command, args []string) {
		baseremote, err := parse.Remote(strings.Split(base, ":")[0])
		if err != nil {
			color.Red("%v", err)
			os.Exit(0)
		}
		urlStr := fmt.Sprintf("https://api.github.com/repos/%s/pulls", baseremote)

		headremote, err := parse.Remote(strings.Split(head, ":")[0])
		if err != nil {
			color.Red("%v", err)
			os.Exit(0)
		}
		userName := strings.Split(headremote, "/")[0]
		head = fmt.Sprintf("%s:%s", userName, strings.Split(head, ":")[1])

		if inBrowser {
			url := fmt.Sprintf("https://github.com/%s/compare/%s...%s?expand=1", baseremote, strings.Split(base, ":")[1], head)
			err := browser.OpenURLInBrowser(url)
			if err != nil {
				color.Red("Failed to open in browser : %+v", err)
				os.Exit(0)
			}

		} else {

			// check for PR title
			if !inEditor && len(args) == 0 {
				color.Red("PR title required")
				os.Exit(0)
			}

			// create a new PR
			newPR := request.PR{Title: args[0]}
			newPR.Base = strings.Split(base, ":")[1]
			newPR.Head = head

			if inEditor {

				var editorOutput []byte
				var err error
				fmt.Println("Waiting for the editor to close....")
				if len(args) != 0 {
					editorOutput, err = editor.OpenEditor(args[0])
				} else {
					editorOutput, err = editor.OpenEditor("")
				}
				if err != nil {
					color.Red("Failed to read from editor : %v", err)
					os.Exit(0)
				}

				// parse the body
				bodyContent := strings.Split(string(editorOutput), "\n\n")[1]
				newPR.Body = bodyContent
			}
			fmt.Println("Creating a PR.....")
			resp, err := request.Request(newPR, urlStr, token)
			if err != nil {
				color.Red("Failed to create a PR : %v", err)
				os.Exit(0)
			}
			resJSON := make(map[string]interface{})
			bytes, _ := ioutil.ReadAll(resp.Body)
			if err := json.Unmarshal(bytes, &resJSON); err != nil {
				color.Red("Failed to parse the response : %v", err)
				os.Exit(0)
			}
			if resp.Status == "201 Created" {
				color.Green("PR created!! :)")
				color.Blue("%s", resJSON["html_url"])
			} else {
				color.Red("Ooops, something went wrong :(")
				color.Red("%s", resJSON["message"])
			}
		}
	},
}

var base string
var head string
var token string
var inEditor bool
var inBrowser bool

func init() {
	cfgFile := fmt.Sprintf("%s/.ghpr.json", os.Getenv("HOME"))
	viper.SetConfigFile(cfgFile)
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	token = viper.GetString("token")
	inEditor = viper.GetBool("inEditor")

	// get the current branch
	currentBranch, err := parse.CurrentBranch()
	if err != nil {
		color.Red("Not a git repositiory")
		currentBranch = "currentbranch"
	}

	// define flags
	f := rootCmd.PersistentFlags()
	f.StringVarP(&base, "base", "B", "upstream:master", "Repo to which the PR is to be made - remotename:branch ")
	f.StringVarP(&head, "head", "H", "origin:"+currentBranch, "Repo in which your changes lie - remotename:branch ")
	f.BoolVarP(&inBrowser, "browser", "b", false, "Open PR creation page in the browser")

}

// Execute executes the command and returns the error
func Execute() error {
	return rootCmd.Execute()
}
