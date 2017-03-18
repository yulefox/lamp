package lamp

import (
	"errors"
	"flag"
	"fmt"

	"github.com/yulefox/lamp/contrib"
	"github.com/yulefox/lamp/internal"

	nsq "github.com/nsqio/go-nsq"
)

var configFilename = "lamp.json"

// Lamp microservice
type Lamp struct {
	config contrib.LampConfig

	// NSQ Consumer
	consumer *nsq.Consumer
}

// On turns on the Lamp
func (l *Lamp) On() error {
	l.config.NSQConfig = nsq.NewConfig()

	if err := l.Reload(); err != nil {
		panic(err)
	}

	c, err := nsq.NewConsumer(l.config.Topic, l.config.Channel, l.config.NSQConfig)
	if err != nil {
		return err
	}
	c.AddHandler(l)
	c.SetLogger(l, nsq.LogLevelInfo)

	if err := c.ConnectToNSQLookupds(l.config.LookupAddrs); err != nil {
		return err
	}
	l.consumer = c
	return nil
}

// Off turns off the Lamp
func (l *Lamp) Off() error {
	l.consumer.Stop()
	return nil
}

// Reload reloads the Lamp
func (l *Lamp) Reload() error {
	if err := internal.LoadConfig(configFilename, &l.config); err != nil {
		return err
	}

	flagSet := flag.NewFlagSet("", flag.ExitOnError)
	flagSet.Var(&nsq.ConfigFlag{Config: l.config.NSQConfig}, "nsq-opt", "option to pass through to nsq.Consumer (may be given multiple times)")

	if err := flagSet.Parse(l.config.NSQConfigFlagSet); err != nil {
		return err
	}
	return nil
}

// Restart restarts the Lamp
func (l *Lamp) Restart() error {
	return errors.New("NO implementation")
}

// Output logger
func (l *Lamp) Output(calldepth int, s string) error {
	fmt.Println(s)
	return nil
}

// HandleMessage message handler for the `gm` topic
func (l Lamp) HandleMessage(msg *nsq.Message) error {
	fmt.Printf("%+v\n", msg)
	return nil
}
