package lamp

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	nsq "github.com/nsqio/go-nsq"
)

type TestLampHandler struct {
}

func (h TestLampHandler) Output(calldepth int, s string) error {
	fmt.Println(s)
	return nil
}

func (h TestLampHandler) HandleMessage(msg *nsq.Message) error {
	fmt.Printf("raw: %v\n", msg)
	if string(msg.Body) == "TOBEFAILED" {
		return errors.New("fail this message")
	}

	data := struct {
		Host string
		Data string
	}{}

	err := json.Unmarshal(msg.Body, &data)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Printf("%+v\n", data)
	return nil
}

func TestLamp(t *testing.T) {
	h := &TestLampHandler{}
	l, err := NewLamp(h)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v\n", l.Consumer.Stats())
	l.On()
	l.Off()

	<-l.Consumer.StopChan
}
