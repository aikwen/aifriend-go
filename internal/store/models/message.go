package models

import "gorm.io/gorm"


type Message struct {
    gorm.Model

    FriendID     uint   `gorm:"not null;index" json:"friend_id"`
    Friend       Friend `gorm:"foreignKey:FriendID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

    UserMessage  string `gorm:"type:text;not null" json:"user_message"`
    Output       string `gorm:"type:text;not null" json:"output"`
}