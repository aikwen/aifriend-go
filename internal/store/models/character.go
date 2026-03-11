package models

import "gorm.io/gorm"


type Character struct {
	gorm.Model
	AuthorID        uint `gorm:"not null;index" json:"author_id"`
	Author          User `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Name            string `gorm:"type:varchar(50);not null" json:"name"`
	Profile         string `gorm:"type:longtext;not null" json:"profile"`
	Photo           string `gorm:"type:varchar(255);not null" json:"photo"`
	BackgroundImage string `gorm:"type:varchar(255);not null" json:"background_image"`
}