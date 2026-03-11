package store

import (
	"gorm.io/gorm"

	"github.com/aikwen/aifriend-go/internal/store/character"
	"github.com/aikwen/aifriend-go/internal/store/friend"
	"github.com/aikwen/aifriend-go/internal/store/user"
)



type Database struct{
	db *gorm.DB
	Character character.Store
	Friend    friend.Store
	User      user.Store
}


func NewDatabase(db *gorm.DB) *Database {
	return &Database{
		db: db,
		Character: character.NewCharacterStore(db),
		Friend: friend.NewFriendStore(db),
		User: user.NewUserStore(db),
	}
}