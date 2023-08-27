package model

import (
	"room-booking/app/database"

	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	Title       string  `gorm:""`
	Address     string  `gorm:""`
	Price       uint64  `gorm:""`
	Area        uint32  `gorm:""`
	IsAvailable bool    `gorm:""`
	Images      []Image `gorm:"many2many:room_images;"`
}

func init() {
	database.DB.AutoMigrate(&Room{})
}
