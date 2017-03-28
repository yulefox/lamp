package lamp

import (
	"testing"

	"os"

	"time"

	nsq "github.com/nsqio/go-nsq"
)

type Hello struct {
}

// EchoServe run echo server
func (s *Hello) EchoServe() {
}

// HandleMessage message handler for the `cmd` topic
func (s *Hello) HandleMessage(msg *nsq.Message) (err error) {
	return
}

func TestLamp(t *testing.T) {
	s := &Hello{}
	go On(s)
	<-time.After(1 * time.Second)
	pid := os.Getpid()
	p, _ := os.FindProcess(pid)
	p.Signal(os.Interrupt)
}
