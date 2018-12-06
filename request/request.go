package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/fatih/color"
)

// PR represents the parameters to be passed to the api as json for creating a pull-request
type PR struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Head  string `json:"head"`
	Base  string `json:"base"`
}

// Request makes a new PR request with the given parameters
func Request(newPR PR, url string, token string) {
	// marshal the newPR
	jsonObj, _ := json.Marshal(&newPR)
	client := &http.Client{}
	r, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonObj)) // URL-encoded payload
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
