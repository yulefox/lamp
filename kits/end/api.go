package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	nsq "github.com/nsqio/go-nsq"
	"github.com/yulefox/lamp"
)

// Serve handle for the `cmd` topic
type Serve struct {
}

// HandleMessage message handler for the `cmd` topic
func (s *Serve) HandleMessage(msg *nsq.Message) (err error) {
	if string(msg.Body) == "TOBEFAILED" {
		return errors.New("fail this message")
	}

	var m url.Values
	err = json.Unmarshal(msg.Body, &m)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		recover()
	}()

	api := m.Get("_api")
	switch api {
	case "get_roles":
		json, err := json.Marshal(m)
		if err != nil {
			return err
		}
		lamp.NSQPublish("localhost:4151", "cmd_res", json)
		break
	default:
		break
	}
	return
}

// EchoServe run echo server
func (s *Serve) EchoServe() {
}

func main() {
	s := &Serve{}
	lamp.On(s)
}
