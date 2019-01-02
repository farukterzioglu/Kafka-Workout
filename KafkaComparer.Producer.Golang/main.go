package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/Shopify/sarama"
)

const (
	broker = "172.31.162.65:9092"
	topic  = "create-review"
)

func main() {
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
	config.ClientID = "tagsProducer"
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	producer, err = sarama.NewSyncProducer([]string{broker}, config)

	return
}

func publish(message string, producer sarama.SyncProducer) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	p, o, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("Error publish: ", err.Error())
	}

	fmt.Printf("Delivered [p:%d] (@%d)\n'", p, o)
}
