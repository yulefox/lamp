package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	nsq "github.com/nsqio/go-nsq"
	"github.com/yulefox/lamp"
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

func main() {
	h := &GMHandler{}
	l, err := lamp.NewLamp(h)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v\n", l.Consumer.Stats())
	l.On()

	<-l.Consumer.StopChan
}
