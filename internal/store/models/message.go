package models

import "gorm.io/gorm"


type Message struct {
    gorm.Model

    FriendID     uint   `gorm:"not null;index" json:"friend_id"`
    Friend       Friend `gorm:"foreignKey:FriendID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

    UserMessage  string `gorm:"type:text;not null" json:"user_message"`
    Input        string `gorm:"type:text;not null" json:"input"`
    Output       string `gorm:"type:text;not null" json:"output"`

    InputTokens  int    `gorm:"default:0" json:"input_tokens"`
    OutputTokens int    `gorm:"default:0" json:"output_tokens"`
    TotalTokens  int    `gorm:"default:0" json:"total_tokens"`
}