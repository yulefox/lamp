package apps

import (
	"testing"
	"time"
)

func TestLampGM(t *testing.T) {
	l := &GM{}
	err := l.On()

	if err != nil {
		return
	}

	<-time.After(60 * time.Second)
	l.Off()
}
