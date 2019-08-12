package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"

	"github.com/sausheong/tanuki/structs"
)

func main() {
	listener, err := net.Listen("tcp", ":0")
	fmt.Println(listener.Addr().String())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	fmt.Printf("Serving %s\n", conn.RemoteAddr().String())
	for {
		data, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		var request structs.RequestInfo
		json.Unmarshal(data, &request)

		response := `{"status": %d, "header": %s, "body": "%s"}`
		headers := make(map[string][]string)
		// this sets a cookie in the header
		headers["Set-Cookie"] = []string{"hello=world; expires=Mon, 12-Dec-2020 20:20:00 GMT"}
		h, _ := json.Marshal(headers)

		resp := fmt.Sprintf(response, 200, h, "hello "+request.Params["name"][0])

		conn.Write([]byte(resp))
	}
}
