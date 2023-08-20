package service

import (
	"context"
	"fmt"
	"room-booking/mapper"
	"room-booking/pb"
	roomRepository "room-booking/repository/room"
	"strings"

	fieldmask_utils "github.com/mennanov/fieldmask-utils"
	"google.golang.org/protobuf/types/known/emptypb"
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
	room := mapper.MapToModel(in)
	id := roomRepository.NewRoom(room).Create()

	fmt.Printf("Room created successfully with id #%d", id)

	return &pb.CreateRoomResponse{
		Id: uint64(id),
	}, nil
}

// FindRoom implements pb.RoomServiceServer
func (s *RoomServer) FindRoom(ctx context.Context, in *pb.FindRoomRequest) (*pb.FindRoomResponse, error) {
	room := roomRepository.FindById(int64(in.Id))

	fmt.Printf("Room with id #%d found successfully", room.ID)

	return &pb.FindRoomResponse{
		Room: mapper.MapToResponse(room),
	}, nil
}

// UpdateRoom implements pb.RoomServiceServer
func (s *RoomServer) UpdateRoom(ctx context.Context, in *pb.UpdateRoomRequest) (*pb.UpdateRoomResponse, error) {
	updatedRoom := in.GetRoom()
	id := updatedRoom.Id
	fm := in.GetFieldMask()

	roomDst := make(map[string]interface{})
	mask, _ := fieldmask_utils.MaskFromProtoFieldMask(fm, func(s string) string {
		return strings.Title(strings.ToLower(s))
	})
	err := fieldmask_utils.StructToMap(mask, updatedRoom, roomDst)

	if err != nil {
		fmt.Println("Error while mapping: ", err)

		return nil, err
	}

	room := roomRepository.Update(int64(id), roomDst)

	fmt.Printf("Room with id #%d has been updated successfully", id)

	return &pb.UpdateRoomResponse{
		Room: mapper.MapToResponse(room),
	}, nil
}

// DeleteRoom implements pb.RoomServiceServer
func (s *RoomServer) DeleteRoom(ctx context.Context, in *pb.DeleteRoomRequest) (*emptypb.Empty, error) {
	id := in.Id
	roomRepository.Delete(int64(id))

	fmt.Printf("Room with id #%d has been deleted successfully", id)

	return &emptypb.Empty{}, nil
}
