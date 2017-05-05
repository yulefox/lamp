package mq

import (
	"log"
	"testing"

	"os"
	"os/signal"

	"github.com/Shopify/sarama"
)

func TestSubscribe(t *testing.T) {
	consumer, err := sarama.NewConsumer([]string{"yulefox.com:9092"}, nil)

	if err != nil {
		panic(err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition("test", 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumed := 0

ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("offset: %d, %s\n", msg.Offset, string(msg.Value))
			consumed++
		case <-signals:
			break ConsumerLoop
		}
	}

	log.Printf("consumed: %d\n", consumed)
}

func TestPublish(t *testing.T) {
	consumer, err := sarama.NewConsumer([]string{"yulefox.com:9092"}, nil)

	if err != nil {
		panic(err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition("test", 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumed := 0

ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("offset: %d, %s\n", msg.Offset, string(msg.Value))
			consumed++
		case <-signals:
			break ConsumerLoop
		}
	}

	log.Printf("consumed: %d\n", consumed)
}
