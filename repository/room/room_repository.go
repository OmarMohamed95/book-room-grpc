package room

import (
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

// Get all rooms
func GetAll(criteria map[string]interface{}) []model.Room {
	var rooms []model.Room
	database.DB.Find(&rooms)

	return rooms
}

// Create a new room
func (r *Room) Create() uint {
	database.DB.Create(&r)

	return r.ID
}

// Find room by id
func FindById(id int64) model.Room {
	var room model.Room
	database.DB.Where("ID = ?", id).Find(&room)

	return room
}

// Update room by id
func Update(id int64, roomUpdates map[string]interface{}) model.Room {
	room := FindById(id)
	database.DB.Model(&room).Where("ID = ?", id).Updates(roomUpdates)

	return room
}

// Delete room by id
func Delete(id int64) {
	var room model.Room
	database.DB.Where("ID = ?", id).Delete(&room)
}
