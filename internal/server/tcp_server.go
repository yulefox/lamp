package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// TCPServer server
type TCPServer struct {
	sync.RWMutex

	// 64bit atomic vars need to be first for proper alignment on 32bit platforms
	clientIDSequence int64
	opts             atomic.Value
	tcpListener      net.Listener
	startTime        time.Time
	exitChan         chan int
	waitGroup        sync.WaitGroup
}

// New new
func New(opts *Options) *TCPServer {
	s := &TCPServer{
		startTime: time.Now(),
		exitChan:  make(chan int),
	}
	s.opts.Store(opts)
	return s
}

func (s *TCPServer) getOpts() *Options {
	return s.opts.Load().(*Options)
}

// GetStartTime GetStartTime
func (s *TCPServer) GetStartTime() time.Time {
	return s.startTime
}

// RealTCPAddr RealTCPAddr
func (s *TCPServer) RealTCPAddr() *net.TCPAddr {
	s.RLock()
	defer s.RUnlock()
	return s.tcpListener.Addr().(*net.TCPAddr)
}

// Serve Serve
func (s *TCPServer) Serve() {
	tcpListener, err := net.Listen("tcp", s.getOpts().TCPAddress)
	if err != nil {
		log.Fatal(fmt.Sprintf("listen (%s) failed - %s", s.getOpts().TCPAddress, err))
		os.Exit(1)
	}
	s.tcpListener = tcpListener
	//go TCPServer(s.tcpListener, s, s.getOpts().Logger)
}

// Exit Exit
func (s *TCPServer) Exit() {
	if s.tcpListener != nil {
		s.tcpListener.Close()
	}

	close(s.exitChan)
}
