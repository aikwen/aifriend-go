package store

import (
	"gorm.io/gorm"

	"github.com/aikwen/aifriend-go/internal/store/cache"
	"github.com/aikwen/aifriend-go/internal/store/character"
	"github.com/aikwen/aifriend-go/internal/store/friend"
	"github.com/aikwen/aifriend-go/internal/store/message"
	"github.com/aikwen/aifriend-go/internal/store/models"
	"github.com/aikwen/aifriend-go/internal/store/systemprompt"
	"github.com/aikwen/aifriend-go/internal/store/user"
)



type Database struct{
	db *gorm.DB
	Character character.Store
	Friend    friend.Store
	User      user.Store
	Message   message.Store
	SystemPrompt systemprompt.Store
	Cache        *cache.Cache
}


func NewDatabase(db *gorm.DB) *Database {
	return &Database{
		db: db,
		Character: character.NewCharacterStore(db),
		Friend: friend.NewFriendStore(db),
		User: user.NewUserStore(db),
		Message: message.NewMessageStore(db),
		SystemPrompt: systemprompt.NewSystemPromptStore(db),
		Cache: cache.New(),
	}
}


// 执行数据库迁移
func RunMigrations(db *gorm.DB) error {
    err := db.AutoMigrate(&models.User{},
		&models.Character{},
		&models.Friend{},
		&models.Message{},
		&models.SystemPrompt{},
		)
	if err != nil {
		return err
	}

	// 全文索引
    if err := characterFullTextIndex(db); err != nil {
        return err
    }
	return nil
}

// 全文索引
func characterFullTextIndex(db *gorm.DB) error {
	var count int64

	checkSQL := `
		SELECT COUNT(1)
		FROM information_schema.statistics
		WHERE table_schema = DATABASE()
		AND table_name = 'characters'
		AND index_name = 'idx_fulltext_search'
		`
	if err := db.Raw(checkSQL).Scan(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	createSQL := `
		ALTER TABLE characters
		ADD FULLTEXT INDEX idx_fulltext_search (name, profile) WITH PARSER ngram
		`
	return db.Exec(createSQL).Error
}