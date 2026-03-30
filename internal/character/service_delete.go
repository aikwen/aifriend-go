package character

import (
	"context"

)

// DeleteCharacter 删除指定的角色
func (s *characterService) DeleteCharacter(ctx context.Context, id uint, authorID uint) error {
	c, err := s.database.Character.GetByIDAndAuthor(ctx, id, authorID)
	if err != nil {
		return err
	}
	// 先删除数据库角色，再删除资源文件
	err = s.database.Character.Delete(ctx, id, authorID)
	if err != nil {
		return err // 如果数据库删除失败，直接返回，保留物理文件
	}

	s.syncer.Enqueue(id)

	if c.Photo != ""{
		_ = s.storage.Delete(c.Photo)
	}

	if c.BackgroundImage != ""{
		_ = s.storage.Delete(c.BackgroundImage)
	}
	return nil
}