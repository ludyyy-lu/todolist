package models

import (
	"gorm.io/gorm"
)

const (
	StatusPending     = "pending"
	StatusInProgress  = "in_progress"
	StatusDone        = "done"
	StatusExpired     = "expired"
)

type Todo struct {
	gorm.Model         //一些共有的字段
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	UserID      uint   `json:"user_id"`
	Tags        []Tag  `gorm:"many2many:todo_tags;" json:"tags"`
}
