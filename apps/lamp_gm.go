package apps

import (
	"encoding/json"
	"errors"
	"fmt"

	nsq "github.com/nsqio/go-nsq"
	"github.com/yulefox/lamp/core"
)

// GMHandler handle for the `gm` topic
type GMHandler struct {
}

// Output logger
func (h GMHandler) Output(calldepth int, s string) error {
	fmt.Println(s)
	return nil
}

// HandleMessage message handler for the `gm` topic
func (h GMHandler) HandleMessage(msg *nsq.Message) error {
	if string(msg.Body) == "TOBEFAILED" {
		return errors.New("fail this message")
	}

	event := struct {
		Event  string `json:"event"`
		ArgA   int32  `json:"arg_a"`
		ArgB   int32  `json:"arg_b"`
		Arg64  int64  `json:"arg_64"`
		ArgStr string `json:"arg_str"`
	}{}

	err := json.Unmarshal(msg.Body, &event)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Printf("%+v\n", event)
	return nil
}

// Run run
func (h GMHandler) Run() error {
	l, err := core.NewLamp(h)

	if err != nil {
		return nil
	}

	l.Off()
	<-l.Consumer.StopChan
	return nil
}
