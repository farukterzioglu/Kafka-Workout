package main

import (
	"log"
	"os"
	"os/signal"

	cluster "github.com/bsm/sarama-cluster"
)

func main() {
	// init (custom) config, enable errors and notifications
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	// Config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// init consumer
	brokers := []string{"172.31.162.65:9092"}
	topics := []string{"tagsPart3"}
	consumer, err := cluster.NewConsumer(brokers, "tags-go-consumers", topics, config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			log.Printf("Error: %s\n", err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			log.Printf("Rebalanced: %+v\n", ntf)
		}
	}()

	var commandEngineService *CommandEngineService
	commandEngineService = NewCommandEngineService()

	// consume messages, watch signals
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				commandEngineService.HandleMessage(msg)
				consumer.MarkOffset(msg, "") // mark message as processed
			}
		case <-signals:
			return
		}
	}
}
