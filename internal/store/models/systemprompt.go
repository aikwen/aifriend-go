package models

import "gorm.io/gorm"


type SystemPrompt struct {
	gorm.Model

	Title        string `gorm:"type:varchar(100);not null" json:"title"`
    OrderNumber  int    `gorm:"default:0;not null" json:"order_number"`
    Prompt       string `gorm:"type:text;not null" json:"prompt"`
}