package lamp

import (
	"testing"
	"time"
)

func TestLamp(t *testing.T) {
	l := &Lamp{}
	err := l.On()

	if err != nil {
		return
	}

	<-time.After(60 * time.Second)
	l.Off()
}
