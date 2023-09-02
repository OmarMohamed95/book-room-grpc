package mapper

import (
	"room-booking/model"
	"room-booking/pb"
)

func MapImagesToResponse(images []model.Image) []*pb.Image {
	var pbis []*pb.Image
	for _, i := range images {
		pbi := &pb.Image{
			Name: i.Name,
			Path: i.Path,
		}

		pbis = append(pbis, pbi)
	}

	return pbis
}
