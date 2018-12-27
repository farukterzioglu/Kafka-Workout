package main

import cluster "github.com/bsm/sarama-cluster"

type IConsumer interface {
	StartConsumer()
	Consume()
}

type AddCmdConsumer struct {
	topicName string
	brokers   []string
}

func NewAddCmdConsumer(_brokers []string) *AddCmdConsumer {
	return &AddCmdConsumer{
		topicName: "addCmd",
		brokers:   _brokers,
	}
}

func (consumer *AddCmdConsumer) StartConsumer() {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true

	_, err := cluster.NewConsumer(consumer.brokers, "", []string{consumer.topicName}, config)
	if err != nil {
		panic(err)
	}
}

func (consumer *AddCmdConsumer) Consume() {

}
