package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/farukterzioglu/KafkaComparer/Review.CommandRpcServer/reviewservice"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 10000, "The server port")
)

func main() {
	serverAddr := fmt.Sprintf("localhost:%d", *port)
	lis, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("Running server at %s...\n", serverAddr)

	grpcServer := grpc.NewServer()
	pb.RegisterReviewServiceServer(grpcServer, NewCommandServer())
	grpcServer.Serve(lis)
}
