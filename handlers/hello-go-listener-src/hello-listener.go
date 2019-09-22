package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"

	structs "github.com/sausheong/tanuki/data"
)

func main() {
	// start a TCP socket server at the given port
	// the port number is the STDIN when starting the listener
	listener, err := net.Listen("tcp", ":"+os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	// listen to connections and send the requests to the handler
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleRequest(conn)
	}
}

// handle the request
func handleRequest(conn net.Conn) {
	response := `{"status": %d, "header": %s, "body": "%s"}` + "\n"
	headers := make(map[string][]string)

	for {
		data, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			fmt.Println("Cannot read from buffer - ", err)
			conn.Write([]byte(fmt.Sprintf(response, 500, "{}", "Cannot read from buffer - "+err.Error())))
			return
		}
		var request structs.RequestInfo
		err = json.Unmarshal([]byte(data), &request)
		if err != nil {
			fmt.Println("Cannot unmarshal JSON - ", err)
			conn.Write([]byte(fmt.Sprintf(response, 500, "{}", "Cannot unmarshal JSON - "+err.Error())))
			return
		}

		// this sets a cookie in the header
		headers["Set-Cookie"] = []string{"hello=world; expires=Mon, 12-Dec-2020 20:20:00 GMT"}
		h, _ := json.Marshal(headers)
		resp := fmt.Sprintf(response, 200, h, "hello "+request.Params["name"][0])
		conn.Write([]byte(resp))
	}
}
