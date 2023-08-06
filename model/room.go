package model

import (
	"room-booking/app/database"

	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	Title       string `gorm:""`
	Address     string `gorm:""`
	Price       int    `gorm:""`
	Area        int    `gorm:""`
	IsAvailable bool   `gorm:""`
}

func init() {
	database.DB.AutoMigrate(&Room{})
}
