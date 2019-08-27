package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"strconv"

	"github.com/sausheong/tanuki/data"
	"github.com/spf13/cobra"
	"github.com/tidwall/pretty"
	"gopkg.in/gookit/color.v1"
)

var listen *int
var filePath *string
var reqFile *string
var prettify *bool

func init() {
	rootCmd.AddCommand(sendCmd)
	listen = sendCmd.Flags().Int("listen", 55771, "listener port")
	filePath = sendCmd.Flags().String("file", "", "path to handler file")
	prettify = sendCmd.Flags().Bool("prettify", false, "prettify the JSON output")
	reqFile = sendCmd.Flags().String("request", "request.json", "path to the JSON request file (can use the Tanuki capture server to create this file).")
}

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a JSON request to a Tanuki handler",
	Long:  `Send a JSON request to a Tanuki handler and receive JSON response. For bins use the file option with the path to the handler file. For listeners use the listen option with the port number to send the request to. For listeners it is assumed you've started the listener separately. This tool is only for bins and local listeners.`,
	Run: func(cmd *cobra.Command, args []string) {
		send()
	},
}

func send() {
	reqJSON, err := ioutil.ReadFile(*reqFile)
	if err != nil {
		fmt.Println("Cannot read request JSON file - ", err)
		return
	}
	fmt.Println("\nREQUEST\n-------")
	fmt.Println(string(reqJSON))

	var output []byte
	if *filePath != "" {
		fmt.Println("Sending request to:", *filePath)
		output, err = exec.Command(*filePath, string(reqJSON)).Output()
		if err != nil {
			color.Red.Println("Cannot execute bin", err)
		}

	} else {
		fmt.Println("sending to listener at port", *listen)
		var conn net.Conn
		conn, err = net.Dial("tcp", ":"+strconv.Itoa(*listen))
		if err != nil {
			fmt.Println("Cannot connect to listener", err)
			return
		}
		fmt.Fprintf(conn, string(reqJSON)+"\n")
		output, err = bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			fmt.Println("Cannot read from listener", err)
			return
		}

	}

	fmt.Println("\nRAW RESPONSE\n------------")
	fmt.Println(string(output))

	if *prettify {
		fmt.Println("\nPRETTIFIED RESPONSE\n-------------------")
		fmt.Println(string(pretty.Pretty(output)))
	}

	var response data.ResponseInfo
	err = json.Unmarshal([]byte(output), &response)
	if err != nil {
		color.Error.Println("ERR: Cannot parse response - ", err)
	} else {
		color.Success.Println("Response parsed successfully.")
	}

}
