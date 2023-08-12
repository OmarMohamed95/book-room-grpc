package main

import (
	"context"
	"flag"
	"log"
	"room-booking/pb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	addr := flag.String("address", "", "the server address")
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connect connect to the server: %v", err)
	}
	defer conn.Close()
	c := pb.NewRoomServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.CreateRoom(ctx, &pb.CreateRoomRequest{Room: &pb.Room{}})
	if err != nil {
		log.Fatalf("err while processing request: %v", err)
	}

	log.Printf("Room created with id #%d", r.GetId())
}
