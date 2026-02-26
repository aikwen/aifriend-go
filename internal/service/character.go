package service

import (
	"context"
	"errors"
	"strings"

	"github.com/aikwen/aifriend-go/internal/models"
	"github.com/aikwen/aifriend-go/internal/store"
	"github.com/aikwen/aifriend-go/pkg/storage"
)

// CreateCharacterParam 包含创建角色所需的全部基础信息。
// Name 和 Profile 为必填项
type CreateCharacterParam struct {
	AuthorID    uint
	Name        string
	Profile     string
	PhotoPath   string
	BgImagePath string
}

// UpdateCharacterParam 包更新角色所需的信息。
// PhotoPath 或 BgImagePath 为空字符串，表示用户没有上传图片。
type UpdateCharacterParam struct {
	ID          uint
	AuthorID    uint   // 必须传入
	Name        string
	Profile     string
	PhotoPath   string // 如果为空，则表示不修改图片
	BgImagePath string // 如果为空，则表示不修改背景
}

// CharacterService 定义角色相关的核心业务逻辑接口。
type CharacterService interface {
	// CreateCharacter 为指定的作者创建一个新的 AI 角色。
	CreateCharacter(ctx context.Context, param *CreateCharacterParam) error

	// GetCharacter 根据角色 ID 获取单个角色的详细信息。
	GetCharacter(ctx context.Context, id uint, authorID uint) (*models.Character, error)

	// UpdateCharacter 更新指定角色的信息。
	UpdateCharacter(ctx context.Context, param *UpdateCharacterParam) error

	// DeleteCharacter 删除指定的角色。
	DeleteCharacter(ctx context.Context, id uint, authorID uint) error
}


// characterService 是 CharacterService 接口的具体实现
type characterService struct {
	store store.CharacterStore
	storage storage.Storage
}

// NewCharacterService 实例化角色服务
func NewCharacterService(store store.CharacterStore, st storage.Storage) CharacterService {
	return &characterService{
		store: store,
		storage: st,
	}
}

// CreateCharacter 为指定的用户创建一个新的 AI 角色
func (s *characterService) CreateCharacter(ctx context.Context, param *CreateCharacterParam) error {
	name := strings.TrimSpace(param.Name)
	profile := strings.TrimSpace(param.Profile)

	if name == "" {
		return errors.New("名字不能为空")
	}
	if profile == "" {
		return errors.New("角色介绍不能为空")
	}

	char := &models.Character{
		AuthorID:        param.AuthorID,
		Name:            name,
		Profile:         profile,
		Photo:           param.PhotoPath,
		BackgroundImage: param.BgImagePath,
	}

	return s.store.Create(ctx, char)
}

// GetCharacter 根据角色 ID 获取单个AI角色的详细信息
func (s *characterService) GetCharacter(ctx context.Context, id uint, authorID uint) (*models.Character, error) {
	return s.store.GetByIDAndAuthor(ctx, id, authorID)
}

// UpdateCharacter 更新指定角色的信息，并处理旧图片的物理删除
func (s *characterService) UpdateCharacter(ctx context.Context, param *UpdateCharacterParam) error {
	char, err := s.store.GetByIDAndAuthor(ctx, param.ID, param.AuthorID)
	if err != nil {
		return err // 找不到或不属于该用户
	}

	name := strings.TrimSpace(param.Name)
	profile := strings.TrimSpace(param.Profile)

	if name == "" {
		return errors.New("名字不能为空")
	}

	if profile == "" {
		return errors.New("角色介绍不能为空")
	}

	char.Name = name
	char.Profile = profile

	// 处理图片更新与旧文件删除
	if param.PhotoPath != "" {
		if char.Photo != "" {
			_ = s.storage.Delete(char.Photo)
		}
		char.Photo = param.PhotoPath
	}

	if param.BgImagePath != "" {
		if char.BackgroundImage != ""{
			_ = s.storage.Delete(char.BackgroundImage)
		}
		char.BackgroundImage = param.BgImagePath
	}

	return s.store.Update(ctx, char)
}

// DeleteCharacter 删除指定的角色
func (s *characterService) DeleteCharacter(ctx context.Context, id uint, authorID uint) error {
	return s.store.Delete(ctx, id, authorID)
}