package store

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/models"
	"gorm.io/gorm"
)

type CharacterStore interface {
	Create(ctx context.Context, character *models.Character) error
	GetByIDAndAuthor(ctx context.Context, id uint, authorID uint) (*models.Character, error)
	Update(ctx context.Context, character *models.Character) error
	Delete(ctx context.Context, id uint, authorID uint) error
}

type characterStore struct {
	db *gorm.DB
}

func NewCharacterStore(db *gorm.DB) CharacterStore {
	return &characterStore{db: db}
}

// Create 创建角色
func (s *characterStore) Create(ctx context.Context, character *models.Character) error {
	return s.db.WithContext(ctx).Create(character).Error
}

// GetByIDAndAuthor 获取单个角色
func (s *characterStore) GetByIDAndAuthor(ctx context.Context, id uint, authorID uint) (*models.Character, error) {
	var c models.Character
	// 获取对应的authorId下的character id
	err := s.db.WithContext(ctx).
		Where("id = ? AND author_id = ?", id, authorID).
		First(&c).Error

	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Update 更新角色信息
func (s *characterStore) Update(ctx context.Context, character *models.Character) error {
	return s.db.WithContext(ctx).Save(character).Error
}

// Delete 删除角色，
func (s *characterStore) Delete(ctx context.Context, id uint, authorID uint) error {
	return s.db.WithContext(ctx).
		Where("id = ? AND author_id = ?", id, authorID).
		Delete(&models.Character{}).Error
}