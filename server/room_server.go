package server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"room-booking/logger"
	"room-booking/mapper"
	"room-booking/pb"
	roomRepository "room-booking/repository/room"
	"room-booking/uploader"
	"room-booking/validator"
	"strings"

	fieldmask_utils "github.com/mennanov/fieldmask-utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	newRoom, err := roomRepository.NewRoom(room).Create()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Room created successfully with id #%d", newRoom.ID)

	return &pb.CreateRoomResponse{
		Id: uint64(newRoom.ID),
	}, nil
}

// FindRoom implements pb.RoomServiceServer
func (s *RoomServer) FindRoom(ctx context.Context, in *pb.FindRoomRequest) (*pb.FindRoomResponse, error) {
	room, err := roomRepository.FindById(uint(in.Id))
	if err != nil {
		return nil, err
	}

	fmt.Printf("Room with id #%d found successfully", room.ID)

	return &pb.FindRoomResponse{
		Room: mapper.MapToResponse(*room),
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

	room, err := roomRepository.Update(uint(id), roomDst)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Room with id #%d has been updated successfully", id)

	return &pb.UpdateRoomResponse{
		Room: mapper.MapToResponse(*room),
	}, nil
}

// DeleteRoom implements pb.RoomServiceServer
func (s *RoomServer) DeleteRoom(ctx context.Context, in *pb.DeleteRoomRequest) (*emptypb.Empty, error) {
	id := in.Id
	err := roomRepository.Delete(uint(id))
	if err != nil {
		return nil, err
	}

	fmt.Printf("Room with id #%d has been deleted successfully", id)

	return &emptypb.Empty{}, nil
}

// UploadImage implements pb.RoomServiceServer
func (s *RoomServer) UploadImage(stream pb.RoomService_UploadImageServer) error {
	req, err := stream.Recv()
	if err != nil {
		err = status.Errorf(codes.Unknown, "cannot receive image info")
		log.Print(err)

		return err
	}

	roomId := req.GetInfo().GetRoomId()
	imageName := req.GetInfo().GetImageName()
	imageType := req.GetInfo().GetImageType()

	err = validator.ValidateType(imageType)
	if err != nil {
		return logger.LogError(status.Error(codes.InvalidArgument, err.Error()))
	}

	image := bytes.Buffer{}
	for {
		err := logger.ContextError(stream.Context())
		if err != nil {
			return err
		}

		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}
		if err != nil {
			return logger.LogError(status.Errorf(codes.Unknown, "cannot receive chunk data: %v", err))
		}

		chunk := req.GetChunkData()

		err = validator.IsSizeExceededMaxSize(chunk, image)
		if err != nil {
			return logger.LogError(status.Error(codes.InvalidArgument, err.Error()))
		}

		_, err = image.Write(chunk)
		if err != nil {
			return logger.LogError(status.Errorf(codes.Internal, "cannot write chunk data: %v", err))
		}
	}

	uploadedImage, err := uploader.HandleUpload(image, imageName, imageType, uint(roomId))
	if err != nil {
		return logger.LogError(status.Errorf(codes.Internal, "cannot write chunk data: %v", err))
	}

	imageSize := image.Len()
	res := &pb.UploadImageResponse{
		Id:   uint64(uploadedImage.ID),
		Size: uint32(imageSize),
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return logger.LogError(status.Errorf(codes.Unknown, "cannot send response: %v", err))
	}

	log.Printf("saved image with id: %d, size: %d, roomId: %d", uploadedImage.ID, imageSize, roomId)

	return nil
}
