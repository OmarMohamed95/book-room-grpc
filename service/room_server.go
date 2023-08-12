package service

import (
	"context"
	"fmt"
	"room-booking/pb"
)

// RoomServer is used to implement pb.RoomServiceServer
type RoomServer struct {
	pb.UnimplementedRoomServiceServer
}

// Return new RoomServer instance
func NewRoomServer() *RoomServer {
	return &RoomServer{}
}

// CreateRoom implements pb.RoomServiceServer
func (s *RoomServer) CreateRoom(ctx context.Context, in *pb.CreateRoomRequest) (*pb.CreateRoomResponse, error) {
	fmt.Println("Room Created Successfully!")

	return &pb.CreateRoomResponse{
		Id: 1,
	}, nil
}
