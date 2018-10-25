package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"

	"github.com/aswinmprabhu/github-pr-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// PR represents the parameters to be passed to the api as json for creating a pull-request
type PR struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Head  string `json:"head"`
	Base  string `json:"base"`
}

func openInBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		return fmt.Errorf("Failed to open the browser : %v", err)
	}
	fmt.Println("Opening in browser.....")
	return nil
}

// rootCmd is the main "ghpr" command
var rootCmd = &cobra.Command{
	Use:   "ghpr [options] <title>",
	Short: "Create github pull requests from the command line",
	Run: func(cmd *cobra.Command, args []string) {

		baseremote, err := utils.ParseRemote(strings.Split(base, ":")[0])
		if err != nil {
			log.Fatal(err)
		}
		urlStr := fmt.Sprintf("https://api.github.com/repos/%s/pulls", baseremote)

		headremote, err := utils.ParseRemote(strings.Split(head, ":")[0])
		if err != nil {
			log.Fatal(err)
		}
		userName := strings.Split(headremote, "/")[0]
		head = fmt.Sprintf("%s:%s", userName, strings.Split(head, ":")[1])

		if inBrowser {
			url := fmt.Sprintf("https://github.com/%s/compare/%s...%s?expand=1", baseremote, strings.Split(base, ":")[1], head)
			err := openInBrowser(url)
			if err != nil {
				log.Fatal(err)
			}

		} else {
			// create a new PR
			newPR := PR{Title: args[0]}
			newPR.Base = strings.Split(base, ":")[1]
			newPR.Head = head

			if inEditor {
				editor := os.Getenv("EDITOR")
				// create a temp file
				tmpDir := os.TempDir()
				tmpFile, tmpFileErr := ioutil.TempFile(tmpDir, "prtitle")
				if tmpFileErr != nil {
					log.Fatalf("Error %s while creating tempFile", tmpFileErr)
				}
				// see if the editor exists
				path, err := exec.LookPath(editor)
				if err != nil {
					log.Fatalf("Error %s while looking for %s\n", path, editor)
				}
				// write the title to the file as the first line
				if len(args) != 0 {
					titleBytes := []byte(args[0])
					if err := ioutil.WriteFile(tmpFile.Name(), titleBytes, 0644); err != nil {
						log.Fatalf("Error while writing to file : %s\n", err)
					}
				}

				cmd := exec.Command(path, tmpFile.Name())
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				// open the file in the editor
				err = cmd.Start()
				if err != nil {
					fmt.Printf("Editor execution failed: %s\n", err)
				}
				fmt.Printf("Waiting for editor to close.....\n")
				err = cmd.Wait()
				if err != nil {
					fmt.Printf("Command finished with error: %v\n", err)
				}
				// read from file
				fileContent, err := ioutil.ReadFile(tmpFile.Name())
				if err != nil {
					fmt.Printf("Error while Reading: %s\n", err)

				}
				// parse the body
				bodyContent := strings.Split(string(fileContent), "\n\n")[1]
				newPR.Body = bodyContent
				if err := os.Remove(tmpFile.Name()); err != nil {
					fmt.Println("Error while deleting the tmp file")

				}
			}
			if !inEditor && len(args) == 0 {
				log.Fatal("PR title required")
			}

			// marshal the newPR
			jsonObj, _ := json.Marshal(&newPR)
			client := &http.Client{}
			r, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer(jsonObj)) // URL-encoded payload
			// set the headers
			AuthVal := fmt.Sprintf("token %s", token)
			r.Header.Add("Authorization", AuthVal)
			r.Header.Add("Content-Type", "application/json")

			// make the req
			resp, err := client.Do(r)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			fmt.Println("Creating a PR.....")
			resJSON := make(map[string]interface{})
			bytes, _ := ioutil.ReadAll(resp.Body)
			if err := json.Unmarshal(bytes, &resJSON); err != nil {
				log.Fatal("Failed to parse the response")
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
	currentBranch, err := utils.CurrentBranch()
	if err != nil {
		log.Fatal(err)
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
