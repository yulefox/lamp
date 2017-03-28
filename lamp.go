package lamp

import (
	"errors"
	"flag"
	"fmt"

	"os"

	"os/signal"

	nsq "github.com/nsqio/go-nsq"
	"github.com/yulefox/lamp/contrib"
)

// Shade interface
type Shade interface {
	EchoServe()
	HandleMessage(message *nsq.Message) error
}

// Lamp microservice
type Lamp struct {
	consumers map[string]*nsq.Consumer
}

// On turns on the Lamp
func On(s Shade) (l *Lamp, err error) {
	l = &Lamp{
		consumers: make(map[string]*nsq.Consumer),
	}

	config := &contrib.LampConfig{}
	err = LoadConfig("lamp.json", config)
	if err != nil {
		panic(err)
	}

	for _, n := range config.Nodes {
		n.NSQConfig = nsq.NewConfig()
		flagSet := flag.NewFlagSet("", flag.ExitOnError)
		flagSet.Var(&nsq.ConfigFlag{Config: n.NSQConfig}, "nsq-opt", "option to pass through to nsq.Consumer (may be given multiple times)")

		err = flagSet.Parse(n.NSQConfigFlagSet)
		if err != nil {
			return
		}

		var c *nsq.Consumer
		c, err = nsq.NewConsumer(n.Topic, n.Channel, n.NSQConfig)
		if err != nil {
			return
		}
		c.AddHandler(s)
		c.SetLogger(l, nsq.LogLevelInfo)

		err = c.ConnectToNSQLookupds(n.LookupAddrs)
		if err != nil {
			return
		}
		l.consumers[n.Topic] = c
	}
	go s.EchoServe()

	c := make(chan os.Signal, 1)
	signal.Notify(c,
		os.Interrupt,
		os.Kill,
	)

	<-c
	return
}

// Restart restarts the Lamp
func (l *Lamp) Restart() (err error) {
	err = errors.New("NO implementation")
	return
}

// Output logger
func (l *Lamp) Output(calldepth int, s string) (err error) {
	fmt.Println(s)
	return
}
