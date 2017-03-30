package server

import (
	"log"
	"net"
	"os"
	"testing"
)

func mustStartServer(opts *Options) (*net.TCPAddr, *TCPServer) {
	opts.TCPAddress = "127.0.0.1:0"
	s := New(opts)
	s.Serve()
	return s.RealTCPAddr(), s
}

func TestStartup(t *testing.T) {
	doneExitChan := make(chan int)
	opts := NewOptions()
	opts.Logger = log.New(os.Stdout, "lamp: ", log.Lshortfile|log.Ltime)
	_, s := mustStartServer(opts)

	exitChan := make(chan int)
	go func() {
		<-exitChan
		s.Exit()
		doneExitChan <- 1
	}()
}
