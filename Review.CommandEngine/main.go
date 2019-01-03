package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/farukterzioglu/KafkaComparer/CommandEngine/Models"
)

func main() {
	var commandEngineService *CommandEngineService
	commandEngineService = NewCommandEngineService()

	// init (custom) config, enable errors and notifications
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	// Config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// init consumer
	brokers := []string{"172.31.162.65:9092"}
	topics := commandEngineService.getTopicList()
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

	msgch := make(chan *sarama.ConsumerMessage)
	go func(channel chan *sarama.ConsumerMessage) {
		var wg sync.WaitGroup
		for newMsg := range channel {
			wg.Add(1)
			go func(msg *sarama.ConsumerMessage) {
				defer wg.Done()

				var (
					ctx    context.Context
					cancel context.CancelFunc
				)
				ctx, cancel = context.WithCancel(context.Background())
				defer cancel()

				request := CommandRequest{
					Msg:        msg,
					ResponseCh: make(chan interface{}),
				}

				commandEngineService.HandleMessage(ctx, request)

				reviewID := models.ReviewIdFromContext(ctx)
				fmt.Printf("Review id : %s\n", reviewID)

				<-request.ResponseCh

				consumer.MarkOffset(msg, "")
			}(newMsg)
		}
		wg.Wait()
	}(msgch)

	// consume messages, watch signals
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				msgch <- msg
			}
		case <-signals:
			close(msgch)
			return
		}
	}
}
