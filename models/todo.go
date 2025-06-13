package models

import (
	"todolist/config"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model  //一些共有的字段
	Title string `json:"title"`
	Description string `json:"description"`
	Status string `json:"status"`
	UserID uint `json:"user_id"` //暂时不写用户模块
}

func AutoMigrate() {
	config.DB.AutoMigrate(&Todo{})
}

func CreateTodo(todo *Todo) error {
	return config.DB.Create(todo).Error
}