package internal

import (
	"log"
	"os"

	"fmt"

	"github.com/takama/daemon"
)

var stdlog, errlog *log.Logger

// Service Service
type Service interface {
	Info() (string, string, []string)
	Serve()
}

func init() {
	stdlog = log.New(os.Stdout, "", 0)
	errlog = log.New(os.Stderr, "", 0)
}

func (s *service) serve() (string, error) {
	prog := os.Args[0]
	usage := fmt.Printf("Usage: %s install | remove | start | stop | status", prog)

	if len(os.Args) > 1 {
		cmd := os.Args[1]
		switch cmd {
		case "install":
			return s.Install()
		case "remove":
			return s.Remove()
		case "start":
			return s.Start()
		case "stop":
			return s.Stop()
		case "status":
			return s.Status()
		default:
			return usage, nil
		}
	}
	return
}

// NewDaemon NewDaemon
func NewDaemon(s Service) {
	name, desc, dep := s.Info()
	srv, err := daemon.New(name, desc, dep...)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}
	status, err := se.serve()
	if err != nil {
		errlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}
