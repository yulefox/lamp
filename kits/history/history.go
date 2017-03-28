package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	nsq "github.com/nsqio/go-nsq"
	"github.com/yulefox/lamp"
)

// History handle for the `History` topic
type History struct {
	lamp.Lamp
}

// HandleMessage message handler for the `History` topic
func (l History) HandleMessage(msg *nsq.Message) error {
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

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, 世界!")
}

func history(c echo.Context) error {
	topic := c.Param("topic")
	content := fmt.Sprintf("Hello, %s", topic)

	return c.String(http.StatusOK, content)
}

func main() {
	l := &History{
		Lamp: lamp.Lamp{
			ConfigFilename: "kits/history/lamp.json",
		},
	}
	err := l.On()

	if err != nil {
		return
	}

	e := echo.New()

	e.GET("/", hello)
	e.GET("/history/:topic", history)
	e.Start(":8000")
}
