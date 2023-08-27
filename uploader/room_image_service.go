package uploader

import (
	"bytes"
	"log"
	"room-booking/model"

	roomRepository "room-booking/repository/room"
)

func HandleUpload(image bytes.Buffer, imageName string, imageType string, roomId uint) (*model.Image, error) {
	path, err := newSQSUploader().upload(image, imageType)
	if err != nil {
		log.Printf("error while uploading iamge to SQS: %s", err)

		return nil, err
	}

	room, err := roomRepository.FindById(uint(roomId))
	if err != nil {
		return nil, err
	}

	updatedImage, err := roomRepository.NewRoom(*room).AddImage(&model.Image{
		Name:  imageName,
		Path:  path,
		Rooms: []model.Room{*room},
	})
	if err != nil {
		return nil, err
	}

	return updatedImage, nil
}
