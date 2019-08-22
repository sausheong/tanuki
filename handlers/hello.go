package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sausheong/tanuki/structs"
)

/*
A sample action bin that processes a request and returns a response
Bins are called by Tanuki for every request and are therefore slower and
less effective, but are also simpler to write.
*/

func main() {
	// the request information are sent by Tanuki through STDIN as JSON
	// the JSON includes all the HTTP 1.1 information including the method
	// headers and body (if any)
	var request structs.RequestInfo

	// your response should be in JSON too and must contain the status, header (can
	// be empty) and body
	var resp string
	response := `{"status": %d, "header": %s, "body": "%s"}`
	headers := make(map[string][]string)

	// unmarshal the request from JSON into a struct
	err := json.Unmarshal([]byte(os.Args[1]), &request)
	if err != nil {
		resp = fmt.Sprintf(response, 500, "{}", err.Error())
		fmt.Print(resp)
		return
	}

	// this sets a cookie in the header
	headers["Set-Cookie"] = []string{"hello=world; expires=Mon, 12-Dec-2020 20:20:00 GMT"}
	h, _ := json.Marshal(headers)

	resp = fmt.Sprintf(response, 200, h, "hello "+request.Params["name"][0])
	// response back to Tanuki is through STDOUT so simply write the JSON
	// back to STDOUT and you're done
	fmt.Print(resp)
}
