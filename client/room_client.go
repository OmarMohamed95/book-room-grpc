package client

import (
	"context"
	"log"
	"reflect"
	"room-booking/pb"
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
		"book":   bookRoom,
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
}

func updateRoom(op OperationParams) {

}

func deleteRoom(op OperationParams) {

}

func bookRoom(op OperationParams) {

}

func uploadRoomImages(op OperationParams) {

}
