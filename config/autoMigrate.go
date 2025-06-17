package config

import "todolist/models"

func AutoMigrate() {
	DB.AutoMigrate(&models.Todo{})
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Tag{})
}