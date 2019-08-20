package cmd

import (
	"testing"
)

func TestGetBins(t *testing.T) {
	bins, err := getBins()
	if err != nil {
		t.Error(err)
	}
	t.Log("test bins:", bins)
	if bins[0].Method != "get" {
		t.Error("Unmarshal failed:", bins[0].Method)
	}
}
