package lamp

import (
	"testing"
	"time"
)

type Hello struct {
	Lamp
	a string
}

func TestLamp(t *testing.T) {
	l := &Hello{Lamp: Lamp{ConfigFilename: "kits/history/lamp.json"}}
	err := l.On()

	if err != nil {
		return
	}

	<-time.After(60 * time.Second)
	l.Off()
}
