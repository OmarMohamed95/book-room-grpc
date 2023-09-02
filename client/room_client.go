package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"room-booking/pb"
	"time"

	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

const (
	imagePath = "images/room.jpg"
)

type OperationParams struct {
	pb.RoomServiceClient
	context.Context
}

// immutable valid room operations map
func getValidOperationsMap() map[string]interface{} {
	return map[string]interface{}{
		"create": createRoom,
		"find":   findRoom,
		"update": updateRoom,
		"delete": deleteRoom,
		"upload": uploadRoomImages,
	}
}

// Perform the requested room operation
func Perform(operation string, op OperationParams) {
	r := reflect.ValueOf(getValidOperationsMap()[operation])

	in := make([]reflect.Value, 1)
	in[0] = reflect.ValueOf(op)

	r.Call(in)
}

func createRoom(op OperationParams) {
	r, err := op.RoomServiceClient.CreateRoom(op.Context, &pb.CreateRoomRequest{Room: &pb.Room{
		Title:       "test room",
		Address:     "test address",
		Price:       200000,
		Area:        200,
		IsAvailable: true,
	}})

	if err != nil {
		log.Fatalf("err while processing request: %v", err)
	}

	log.Printf("Room created with id #%d", r.GetId())
}

func findRoom(op OperationParams) {
	roomId := 1
	r, err := op.RoomServiceClient.FindRoom(op.Context, &pb.FindRoomRequest{Id: uint64(roomId)})

	if err != nil {
		log.Fatalf("err while processing request: %v", err)
	}

	log.Printf("Room with id #%d found successfully", roomId)
	log.Printf("Title: %s", r.Room.Title)
	log.Printf("Address: %s", r.Room.Address)
	log.Printf("Price: %d", r.Room.Price)
	log.Printf("Area: %d", r.Room.Area)
	log.Printf("IsAvailable: %t", r.Room.IsAvailable)
	for _, image := range r.Room.Images {
		log.Printf("Image: %s", image)
	}
}

func updateRoom(op OperationParams) {
	fm, err := fieldmaskpb.New(&pb.Room{}, "title")
	if err != nil {
		fmt.Printf("Error while appending paths: %s", err)
	}

	roomId := 1
	op.RoomServiceClient.UpdateRoom(op.Context, &pb.UpdateRoomRequest{
		Room: &pb.Room{
			Id:          uint64(roomId),
			Title:       "update test room title",
			Address:     "test address",
			Price:       200000,
			Area:        200,
			IsAvailable: true,
		},
		FieldMask: fm,
	})

	fmt.Printf("Room with id #%d has been updated successfully", roomId)
}

func deleteRoom(op OperationParams) {
	roomId := 1
	_, err := op.RoomServiceClient.DeleteRoom(op.Context, &pb.DeleteRoomRequest{Id: uint64(roomId)})

	if err != nil {
		log.Fatalf("err while processing request: %v", err)
	}

	fmt.Printf("Room with id #%d has been deleted successfully", roomId)
}

func uploadRoomImages(op OperationParams) {
	roomId := 1
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("cannot open image file: ", err)
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := op.RoomServiceClient.UploadImage(ctx)
	if err != nil {
		log.Fatal("cannot upload image: ", err)
	}

	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Info{
			Info: &pb.ImageInfo{
				RoomId:    uint64(roomId),
				ImageName: filepath.Base(imagePath),
				ImageType: filepath.Ext(imagePath),
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		log.Fatal("cannot send image info to server: ", err, stream.RecvMsg(nil))
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("cannot read chunk to buffer: ", err)
		}

		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			log.Fatal("cannot send chunk to server: ", err, stream.RecvMsg(nil))
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("cannot receive response: ", err)
	}

	log.Printf("Image uploaded with id: %d, size: %d", res.GetId(), res.GetSize())
}
