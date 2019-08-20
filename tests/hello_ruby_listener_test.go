package tests

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestGetHelloRubyListener(t *testing.T) {
	port := 59559
	go exec.Command("../listeners/hello-ruby", strconv.Itoa(port)).Run()
	time.Sleep(time.Second)
	conn, err := net.Dial("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		t.Error("Cannot connect to listener", err)
	}
	fmt.Fprintf(conn, string(reqJSON)+"\n")
	// listen for reply
	output, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		t.Error("Cannot read from listener", err)
	}
	if strings.Compare(string(output), expected) != 0 {
		t.Error("Output is not what is expected - ", string(output))
	}
}
