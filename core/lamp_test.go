package core

import (
	"fmt"
	"testing"

	nsq "github.com/nsqio/go-nsq"
)

type MyHandler struct {
}

// Output logger
func (h MyHandler) Output(calldepth int, s string) error {
	fmt.Println(s)
	return nil
}

// HandleMessage message handler for the `gm` topic
func (h MyHandler) HandleMessage(msg *nsq.Message) error {
	return nil
}

func (h MyHandler) Run() error {
	l, err := NewLamp(h)

	if err != nil {
		return nil
	}

	l.Off()
	<-l.Consumer.StopChan
	return nil
}

func TestLampRun(t *testing.T) {
	h := MyHandler{}
	h.Run()
}
