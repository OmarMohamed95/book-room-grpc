package mapper

import (
	"room-booking/model"
	"room-booking/pb"
)

// Map room request to room model
func Map(rr *pb.CreateRoomRequest) model.Room {
	return model.Room{
		Title:       rr.GetRoom().GetTitle(),
		Address:     rr.GetRoom().GetAddress(),
		Price:       rr.GetRoom().GetPrice(),
		Area:        rr.GetRoom().GetArea(),
		IsAvailable: rr.GetRoom().GetIsAvailable(),
	}
}
