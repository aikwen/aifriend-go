package character

import (
	"context"

)

// DeleteCharacter 删除指定的角色
func (s *characterService) DeleteCharacter(ctx context.Context, id uint, authorID uint) error {
	c, err := s.characterStore.getByIDAndAuthor(ctx, id, authorID)
	if err != nil {
		return err
	}
	// 先删除数据库角色，再删除资源文件
	err = s.characterStore.delete(ctx, id, authorID)
	if err != nil {
		return err // 如果数据库删除失败，直接返回，保留物理文件
	}

	if c.Photo != ""{
		_ = s.storage.Delete(c.Photo)
	}

	if c.BackgroundImage != ""{
		_ = s.storage.Delete(c.BackgroundImage)
	}
	return nil
}