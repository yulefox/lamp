package core

import (
	"errors"
	"flag"

	"github.com/yulefox/lamp/contrib"

	nsq "github.com/nsqio/go-nsq"
)

var configFilename = "lamp.json"

// Lamp microservice
type Lamp struct {
	config contrib.LampConfig

	// NSQ Consumer
	Consumer *nsq.Consumer
}

// ILamp Lamp interface
type ILamp interface {
	Run() error
	Output(calldepth int, s string) error
	HandleMessage(message *nsq.Message) error
}

// NewLamp creates a new instance of Lamp as a microservice
//
//  Lamp will be turned on while created successfully.
func NewLamp(i ILamp) (*Lamp, error) {
	l := new(Lamp)
	l.config.NSQConfig = nsq.NewConfig()

	if err := l.Reload(); err != nil {
		panic(err)
	}

	c, err := nsq.NewConsumer(l.config.Topic, l.config.Channel, l.config.NSQConfig)
	if err != nil {
		return nil, err
	}
	c.AddHandler(i)
	c.SetLogger(i, nsq.LogLevelInfo)
	l.Consumer = c

	if err := l.On(); err != nil {
		return nil, err
	}
	return l, nil
}

// On turns on the Lamp
func (l *Lamp) On() error {
	if err := l.Consumer.ConnectToNSQLookupds(l.config.LookupAddrs); err != nil {
		return err
	}
	return nil
}

// Off turns off the Lamp
func (l *Lamp) Off() error {
	l.Consumer.Stop()
	return nil
}

// Reload reloads the Lamp
func (l *Lamp) Reload() error {
	if err := LoadConfig(configFilename, &l.config); err != nil {
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
