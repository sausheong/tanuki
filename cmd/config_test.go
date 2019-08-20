package cmd

import (
	"testing"
)

func TestGetListeners(t *testing.T) {
	handlers, err := getHandlers("../handlers.yaml")
	if err != nil {
		t.Error(err)
	}
	t.Log("handlers:", handlers)
	if handlers[0].Method != "get" {
		t.Error("Unmarshal failed:", handlers[0].Method)
	}
}
