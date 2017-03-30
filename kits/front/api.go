package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	nsq "github.com/nsqio/go-nsq"
	"github.com/yulefox/lamp"
)

// Req request message
type Req struct {
	Token string
	Res   chan []byte
}

var (
	reqC chan Req
	resC chan url.Values
)

func init() {
	reqC = make(chan Req, 1024)
	resC = make(chan url.Values, 1024)
}

func serveAPIv1(c echo.Context) (err error) {
	topic := c.Param("topic")
	api := c.Param("api")
	params := c.QueryParams()

	req := Req{
		Token: uuid.New().String(),
		Res:   make(chan []byte),
	}
	params.Set("_api", api)
	params.Set("_token", req.Token)
	json, err := json.Marshal(params)

	if err != nil {
		return
	}
	lamp.NSQPublish("localhost:4151", topic, json)
	reqC <- req
	res, ok := <-req.Res
	if ok {
		return c.JSON(http.StatusOK, (string)(res))
	}
	return c.JSON(http.StatusBadRequest, "invalid request")
}

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

	resC <- m
	return
}

// EchoServe run echo server
func (s *Serve) EchoServe() {
	e := echo.New()
	g := e.Group("/api/v1")

	g.GET("/:topic/:api", serveAPIv1)
	e.Start(":8000")
}

func serve() {
	reqs := make(map[string]Req)

	for {
		select {
		case req, ok := <-reqC:
			if ok {
				reqs[req.Token] = req
			}
			break
		case res, ok := <-resC:
			if ok {
				token := res.Get("_token")
				res.Del("_token")
				res.Del("_api")
				req, ok := reqs[token]
				if ok {
					json, err := json.Marshal(res)
					if err != nil {
						req.Res <- nil
					}
					req.Res <- json
				}
			}
			break
		}
	}
}

func main() {
	go serve()
	s := &Serve{}
	lamp.On(s)
}
