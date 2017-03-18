package apps

import (
	"testing"
	"time"

	"github.com/yulefox/lamp/core"
)

func TestLampGM(t *testing.T) {
	h := &GMHandler{}
	l, err := core.NewLamp(h)

	if err != nil {
		return
	}

	<-time.After(10 * time.Second)
	l.Off()
}
