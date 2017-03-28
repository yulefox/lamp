package contrib

import nsq "github.com/nsqio/go-nsq"

// LampNode for Lamp
type LampNode struct {
	Topic            string   `json:"topic"`
	Channel          string   `json:"channel"`
	NSQDAddrs        []string `json:"nsqd_addrs"`
	LookupAddrs      []string `json:"nsqlookupd_addrs"`
	LogLevel         string   `json:"log_level"`
	NSQConfigFlagSet []string `json:"nsq_config"`
	NSQConfig        *nsq.Config
}

// LampConfig for Lamp
type LampConfig struct {
	Nodes []LampNode `json:"lamp"`
}
