package models

import "todolist/config"

func AutoMigrate() {
	config.DB.AutoMigrate(&Todo{})
	config.DB.AutoMigrate(&User{})
}