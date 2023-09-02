package uploader

import (
	"bytes"
	"log"
	"room-booking/model"
	"room-booking/pb"
	"room-booking/validator"

	roomRepository "room-booking/repository/room"
)

func HandleUpload(image bytes.Buffer, imageInfo *pb.ImageInfo, roomId uint, uploader ImageUploader) (*model.Image, error) {
	err := validator.ValidateType(imageInfo.GetImageType())
	if err != nil {
		return nil, err
	}

	path, err := uploader.upload(image, imageInfo)
	if err != nil {
		log.Printf("error while uploading iamge to S3: %s", err)

		return nil, err
	}

	room, err := roomRepository.FindById(uint(roomId))
	if err != nil {
		return nil, err
	}

	updatedImage, err := roomRepository.NewRoom(*room).AddImage(&model.Image{
		Name:  imageInfo.GetImageName(),
		Path:  path,
		Rooms: []model.Room{*room},
	})
	if err != nil {
		return nil, err
	}

	return updatedImage, nil
}
