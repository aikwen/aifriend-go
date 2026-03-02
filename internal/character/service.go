package character

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/models"
	"github.com/aikwen/aifriend-go/pkg/storage"
	"gorm.io/gorm"
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
	AuthorID    uint // 必须传入
	Name        string
	Profile     string
	PhotoPath   string // 如果为空，则表示不修改图片
	BgImagePath string // 如果为空，则表示不修改背景
}

// Service 定义角色相关的核心业务逻辑接口。
type Service interface {
	// CreateCharacter 为指定的作者创建一个新的 AI 角色。
	CreateCharacter(ctx context.Context, param *CreateCharacterParam) error

	// GetCharacter 根据角色 ID 获取单个角色的详细信息。
	GetCharacter(ctx context.Context, id uint, authorID uint) (*models.Character, error)

	// UpdateCharacter 更新指定角色的信息。
	UpdateCharacter(ctx context.Context, param *UpdateCharacterParam) error

	// DeleteCharacter 删除指定的角色。
	DeleteCharacter(ctx context.Context, id uint, authorID uint) error

	// GetCharacterList 获取指定数量的虚拟角色
	GetCharacterList(ctx context.Context, authorID uint, offset int) ([]*models.Character, error)

	// 根据关键字搜索角色
	SearchCharacters(ctx context.Context, offset int, limit int, searchQuery string) ([]*models.Character, error)

	Exist(ctx context.Context, characterId uint) (bool, error)
}

// characterService 是 CharacterService 接口的具体实现
type characterService struct {
	characterStore store
	storage        storage.Storage
}

// NewCharacterService 实例化角色服务
func NewCharacterService(db *gorm.DB, st storage.Storage) Service {
	return &characterService{
		characterStore: newCharacterStore(db),
		storage:        st,
	}
}
