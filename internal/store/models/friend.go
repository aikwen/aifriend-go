package models


import "gorm.io/gorm"


type Friend struct {
	gorm.Model

	UserID      uint      `gorm:"not null;index" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	CharacterID uint      `gorm:"not null;index" json:"character_id"`
	Character   Character `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	Memory      string    `gorm:"type:text" json:"memory"`
}