package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/farukterzioglu/KafkaComparer/Review.CommandEngine/Models"
	pb "github.com/farukterzioglu/KafkaComparer/Review.CommandRpcServer/reviewservice"
	"google.golang.org/grpc"
)

var (
	serverAddr   = flag.String("server_addr", "127.0.0.1:10000", "The rpc server address in the format of host:port")
	kafkaBrokers = flag.String("kafka_brokers", "127.0.0.1:9092", "The kafka broker address in the format of host:port")
)

func main() {
	flag.Parse()

	// Configure gRpc client
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
		panic(err)
	}
	defer conn.Close()

	client := pb.NewReviewServiceClient(conn)

	// Configure command engine service
	var commandEngineService *CommandEngineService
	commandEngineService = NewCommandEngineService(&client)

	// init (custom) config, enable errors and notifications
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	// Config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// init kafka consumer
	brokers := []string{*kafkaBrokers} // TODO : parse comma seperated broker list
	topics := commandEngineService.getTopicList()
	consumer, err := cluster.NewConsumer(brokers, "review-command-engine", topics, config)
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

	failedMsgChn := make(chan *sarama.ConsumerMessage)
	go func(failedMsgCh chan *sarama.ConsumerMessage) {
		for msg := range failedMsgChn {
			fmt.Printf("Message failed : %s\n", msg.Value)
			// TODO : Process failed messages again
		}
	}(failedMsgChn)

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

				select {
				case resp := <-request.ResponseCh:
					returnValue := resp.(string)
					fmt.Printf("Review id : %s\n", returnValue)

					reviewID := models.ReviewIdFromContext(ctx)
					fmt.Printf("Review id from context: %s\n", reviewID)
				case err := <-request.ErrCh:
					fmt.Printf("Request failed : %s\n", err.Error())
					failedMsgChn <- msg
				case <-time.After(time.Minute):
					fmt.Printf("Request timedout!\n")
					failedMsgChn <- msg
					cancel()
				}

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
