package main

import "testing"

func TestFix(t *testing.T) {
	servers := []string{
		"7002",
		"7003",
	}

	args := []string{"fix", "db.json"}

	for _, s := range servers {
		args[2] = s
		go Run(args)
	}
}
