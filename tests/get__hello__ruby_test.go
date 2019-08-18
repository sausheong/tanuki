package tests

import (
	"os/exec"
	"strings"
	"testing"
)

func TestGetHelloRubyBin(t *testing.T) {
	expected := `{"status":200,"header":{},"body":"hello sausheong"}` + "\n"
	output, err := exec.Command("../bin/get__hello__ruby", reqJSON).Output()
	if err != nil {
		t.Error("Cannot execute binary", err)
	}
	if strings.Compare(string(output), expected) != 0 {
		t.Error("Output is not what is expected - ", string(output))
	}
}
