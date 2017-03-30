package server

import (
	"log"
)

// Options Options
type Options struct {
	TCPAddress string `flag:"tcp-address"`
	Logger     *log.Logger
}

// NewOptions NewOptions
func NewOptions() *Options {
	return &Options{
		TCPAddress: "0.0.0.0:4150",
	}
}
