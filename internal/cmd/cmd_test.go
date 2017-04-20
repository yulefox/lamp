package cmd

import (
	"testing"
)

func TestLS(t *testing.T) {
	LS(".")
}

func TestSCP(t *testing.T) {
	SCP("20202", "test.json", "juice@sdk.91juice.com:/juice/gmproxy/json")
}
