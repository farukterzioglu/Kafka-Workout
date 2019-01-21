package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	cluster "github.com/bsm/sarama-cluster"
)

var (
	topicName    = flag.String("topic_name", "", "Name of topic to publish")
	kafkaBrokers = flag.String("kafka_brokers", "172.24.96.1:9092", "The kafka broker address in the format of host:port")
)

func main() {
	flag.Parse()
	fmt.Printf("Broker address : %s\n", *kafkaBrokers)
	fmt.Printf("Topic name : %s\n", *topicName)

	// init (custom) config, enable errors and notifications
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	// Config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// init consumer
	// TODO : Support multiple broker list
	brokers := []string{*kafkaBrokers}
	topics := []string{*topicName}
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

	// consume messages, watch signals
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				consumer.MarkOffset(msg, "") // mark message as processed
			}
		case <-signals:
			return
		}
	}
}
