package core

import (
	"testing"
	"time"
)

type MyHandler struct {
	Delegate
}

func TestLampRun(t *testing.T) {
	h := &MyHandler{}
	l, err := NewLamp(h)

	if err != nil {
		return
	}

	<-time.After(10 * time.Second)
	l.Off()
}
