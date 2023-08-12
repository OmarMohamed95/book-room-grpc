package service

import (
	"context"
	"fmt"
	"room-booking/mapper"
	"room-booking/pb"
	roomRepository "room-booking/repository/room"
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
	room := mapper.Map(in)
	id := roomRepository.NewRoom(room).Create()

	fmt.Printf("Room created successfully with id #%d", id)

	return &pb.CreateRoomResponse{
		Id: uint64(id),
	}, nil
}
