package room

import (
	"fmt"
	"room-booking/app/database"
	"room-booking/model"
)

type Room struct {
	model.Room
}

// Return new Room instance
func NewRoom(r model.Room) *Room {
	return &Room{
		Room: r,
	}
}

// Create a new room
func (r *Room) Create() (*Room, error) {
	db := database.DB.Create(&r.Room)
	if db.Error != nil {
		fmt.Printf("error while creating new room: %s", db.Error)
		return nil, db.Error
	}

	return r, nil
}

// Find room by id
func FindById(id uint) (*model.Room, error) {
	var room model.Room
	db := database.DB.Where("ID = ?", id).First(&room)
	if db.Error != nil {
		fmt.Printf("error while finding room with id #%d: %s", id, db.Error)
		return nil, db.Error
	}

	return &room, nil
}

// Update room by id
func Update(id uint, roomUpdates map[string]interface{}) (*model.Room, error) {
	room, err := FindById(id)
	if err != nil {
		return nil, err
	}

	db := database.DB.Model(&room).Where("ID = ?", id).Updates(roomUpdates)
	if db.Error != nil {
		fmt.Printf("error while updating room with id #%d: %s", id, db.Error)
		return nil, db.Error
	}

	return room, nil
}

// Delete room by id
func Delete(id uint) error {
	var room model.Room
	db := database.DB.Where("ID = ?", id).Delete(&room)
	if db.Error != nil {
		fmt.Printf("error while deleting room with id #%d: %s", id, db.Error)
		return db.Error
	}

	return nil
}

// Add new image room
func (room *Room) AddImage(image *model.Image) (*model.Image, error) {
	err := database.DB.Model(&room).Association("Images").Append(image)
	if err != nil {
		fmt.Printf("error while adding a new room image: %s", err)

		return nil, err
	}

	return image, nil
}
