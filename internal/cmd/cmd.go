package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/labstack/gommon/log"
)

func cmd(c, args string) {
	out, err := exec.Command(c, strings.Split(args, " ")...).Output()

	if err != nil {
		log.Error(err)
	}
	fmt.Println(string(out))
}

// LS .
func LS(path string) {
	cmd("ls", path)
}

// SCP .
func SCP(port, fileName, host string) {
	args := fmt.Sprintf("-P %s %s %s", port, fileName, host)
	cmd("scp", args)
}
