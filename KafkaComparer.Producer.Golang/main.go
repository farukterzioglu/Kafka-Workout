package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Shopify/sarama"
)

var (
	topicName    = flag.String("topic_name", "", "Name of topic to publish")
	kafkaBrokers = flag.String("kafka_brokers", "localhost:9092", "The kafka broker address in the format of host:port")
)

func main() {
	flag.Parse()
	fmt.Printf("Broker address : %s\n", *kafkaBrokers)
	fmt.Printf("Topic name : %s\n", *topicName)

	producer, err := initProducer()
	if err != nil {
		fmt.Println("Error while creating producer : ", err.Error())
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter message: ")
		msg, _ := reader.ReadString('\n')

		publish(msg, producer)
	}
}

func initProducer() (producer sarama.SyncProducer, err error) {
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)

	config := sarama.NewConfig()
	config.ClientID = "SampleProducer"
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	producer, err = sarama.NewSyncProducer([]string{*kafkaBrokers}, config)

	return
}

func publish(message string, producer sarama.SyncProducer) {
	msg := &sarama.ProducerMessage{
		Topic: *topicName,
		Value: sarama.StringEncoder(message),
	}

	p, o, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("Error publish: ", err.Error())
	}

	fmt.Printf("Delivered [p:%d] (@%d)\n'", p, o)
}
