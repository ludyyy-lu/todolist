package models

import (
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model         //一些共有的字段
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	UserID      uint   `json:"user_id"`
}
