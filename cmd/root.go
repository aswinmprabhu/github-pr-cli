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
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type PR struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Head  string `json:"head"`
	Base  string `json:"base"`
}

// rootCmd is the main "ghpr" command
var rootCmd = &cobra.Command{
	Use:   "ghpr [OPTIONS] <title>",
	Short: "GitHub create a PR tool for command line",
	Run: func(cmd *cobra.Command, args []string) {
		// create a new PR
		newPR := PR{Title: args[0], Base: "master"}
		if inEditor {
			editor := os.Getenv("EDITOR")
			fmt.Println(editor)
			// create a temp file
			tmpDir := os.TempDir()
			tmpFile, tmpFileErr := ioutil.TempFile(tmpDir, "prtitle")
			if tmpFileErr != nil {
				fmt.Printf("Error %s while creating tempFile", tmpFileErr)
			}
			// see if the editor exists
			path, err := exec.LookPath(editor)
			if err != nil {
				fmt.Printf("Error %s while looking up for %s\n", path, editor)
			}
			// write the title to the file as the first line
			if len(args) != 0 {
				titleBytes := []byte(args[0])
				if err := ioutil.WriteFile(tmpFile.Name(), titleBytes, 0644); err != nil {
					fmt.Printf("Error while writing to file : %s\n", err)
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
			if Debug {
				fmt.Println("Body:", bodyContent)
			}
			newPR.Body = bodyContent
			if err := os.Remove(tmpFile.Name()); err != nil {
				fmt.Println("Error while deleting the tmp file")

			}
		}
		if !inEditor && len(args) == 0 {
			fmt.Println("PR title required")
			os.Exit(0)
		}
		if Debug {
			fmt.Println("Remote:", Remote, "Branch:", Branch, "Title:", args[0])
		}
		// exec "git remote -v" to get the remotes
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
		// parse the repo as username/reponame
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
			newPR.Head = head
		} else {
			newPR.Head = Branch
		}
		// marshal the newPR
		jsonObj, _ := json.Marshal(&newPR)
		client := &http.Client{}
		r, _ := http.NewRequest("POST", urlStr, bytes.NewBuffer(jsonObj)) // URL-encoded payload
		// set the headers
		AuthVal := fmt.Sprintf("token %s", Token)
		r.Header.Add("Authorization", AuthVal)
		r.Header.Add("Content-Type", "application/json")

		// make the req
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
var inEditor bool

func init() {
	cfgFile := fmt.Sprintf("%s/.ghpr.json", os.Getenv("HOME"))
	viper.SetConfigFile(cfgFile)
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	Debug = viper.GetBool("debug")
	Token = viper.GetString("token")
	inEditor = viper.GetBool("inEditor")
	if Debug {
		fmt.Println("Debug:", Debug, "Token:", Token, "inEditor:", inEditor)
	}
	// define flags
	f := rootCmd.PersistentFlags()
	f.StringVarP(&Remote, "remote", "r", "upstream", "Remote GitHub repo to which the PR is to be made")
	f.StringVarP(&Branch, "branch", "b", "master", "The branch from which the PR is to be made")
}

// Execute executes the command and returns the error
func Execute() error {
	return rootCmd.Execute()
}
