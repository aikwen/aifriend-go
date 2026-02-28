package character

import (
	"context"
	"errors"
	"strings"

	"github.com/aikwen/aifriend-go/internal/models"
)

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

	return s.characterStore.create(ctx, char)
}