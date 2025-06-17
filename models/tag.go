package models

//要不要加json字段啊

type Tag struct {
	ID    uint   `gorm:"primaryKey" json:"id"` // 主键ID，自增
	Name  string `gorm:"unique;not null" json:"name"` // 标签名称，唯一且不能为空
	UserID uint  `gorm:"not null" json:"user_id"`
}