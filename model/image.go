package model

import (
	"room-booking/app/database"

	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	Name  string `gorm:""`
	Path  string `gorm:""`
	Rooms []Room `gorm:"many2many:room_images;"`
}

func init() {
	database.DB.AutoMigrate(&Image{})
}
