package mapper

import (
	"room-booking/model"
	"room-booking/pb"
)

// Map room request to room model
func MapToModel(rr *pb.CreateRoomRequest) model.Room {
	return model.Room{
		Title:       rr.GetRoom().GetTitle(),
		Address:     rr.GetRoom().GetAddress(),
		Price:       rr.GetRoom().GetPrice(),
		Area:        rr.GetRoom().GetArea(),
		IsAvailable: rr.GetRoom().GetIsAvailable(),
	}
}

// Map room model to room response
func MapToResponse(rm model.Room) *pb.Room {
	return &pb.Room{
		Title:       rm.Title,
		Address:     rm.Address,
		Price:       rm.Price,
		Area:        rm.Area,
		IsAvailable: rm.IsAvailable,
	}
}
