package models


import (
	"gorm.io/gorm"
)

type User struct {
	// 自动包含 ID, CreatedAt, UpdatedAt, DeletedAt 四个字段
	gorm.Model

	Username string `gorm:"type:varchar(150);uniqueIndex;not null"`
	Password string `gorm:"type:varchar(128);not null"`

	Photo   string `gorm:"type:varchar(255);default:'user/photos/default.png'"`
	Profile string `gorm:"type:varchar(500);default:'该用户什么都没留下'"`
}