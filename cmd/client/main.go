package main

import (
	"context"
	"flag"
	"log"
	"room-booking/client"
	"room-booking/pb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	addr := flag.String("address", "", "the server address")
	operation := flag.String("operation", "", "the server address")
	flag.Parse()

	client.Validate(*operation)

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connect connect to the server: %v", err)
	}
	defer conn.Close()
	c := pb.NewRoomServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client.Perform(*operation, client.OperationParams{
		RoomServiceClient: c,
		Context:           ctx,
	})
}
