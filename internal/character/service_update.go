package character

import (
	"context"
	"errors"
	"strings"

)

// UpdateCharacter 更新指定角色的信息，并处理旧图片的物理删除
func (s *characterService) UpdateCharacter(ctx context.Context, param *UpdateCharacterParam) error {
	char, err := s.characterStore.getByIDAndAuthor(ctx, param.ID, param.AuthorID)
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

	return s.characterStore.update(ctx, char)
}