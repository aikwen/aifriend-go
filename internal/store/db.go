package store

import (
	"gorm.io/gorm"

	"github.com/aikwen/aifriend-go/internal/store/character"
	"github.com/aikwen/aifriend-go/internal/store/friend"
	"github.com/aikwen/aifriend-go/internal/store/user"
	"github.com/aikwen/aifriend-go/internal/store/message"
	"github.com/aikwen/aifriend-go/internal/store/systemprompt"
)



type Database struct{
	db *gorm.DB
	Character character.Store
	Friend    friend.Store
	User      user.Store
	Message   message.Store
	SystemPrompt systemprompt.Store
}


func NewDatabase(db *gorm.DB) *Database {
	return &Database{
		db: db,
		Character: character.NewCharacterStore(db),
		Friend: friend.NewFriendStore(db),
		User: user.NewUserStore(db),
		Message: message.NewMessageStore(db),
		SystemPrompt: systemprompt.NewSystemPromptStore(db),
	}
}