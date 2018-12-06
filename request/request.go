package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// PR represents the parameters to be passed to the api as json for creating a pull-request
type PR struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Head  string `json:"head"`
	Base  string `json:"base"`
}

// Request makes a new PR request with the given parameters
func Request(newPR PR, url string, token string) (*http.Response, error) {
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
		return new(http.Response), fmt.Errorf("Couldn't make the request : %v", err)
	}
	defer resp.Body.Close()
	return resp, nil
}
