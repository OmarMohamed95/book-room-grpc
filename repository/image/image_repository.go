package image

import (
	"fmt"
	"room-booking/app/database"
	"room-booking/model"
)

type Image struct {
	model.Image
}

// Return new Image instance
func NewImage(i model.Image) *Image {
	return &Image{
		Image: i,
	}
}

// Find room images
func FindImagesByRoomId(id uint) (*[]model.Image, error) {
	var images []model.Image
	db := database.DB.Table(
		"images",
	).Select(
		"images.name, images.path",
	).Joins(
		"left join room_images on room_images.image_id = images.id",
	).Where(
		"room_images.room_id = ?", id,
	).Scan(&images)
	if db.Error != nil {
		fmt.Printf("error while finding images for room with id #%d: %s", id, db.Error)
		return nil, db.Error
	}

	return &images, nil
}
