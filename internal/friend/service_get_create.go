package friend

import (
	"context"
	"errors"
)


func (s *friendService) GetOrCreate(ctx context.Context, userID uint, characterID uint) (*FriendDTO, error) {

	if userID == 0 || characterID == 0 {
		return nil, errors.New("无效的用户或角色参数")
	}

	// 查询虚拟角色是否存在
	exists, err := s.charProvider.Exist(ctx, characterID)
	if err != nil {
		return nil, err // 数据库查询异常
	}
	if !exists {
		return nil, errors.New("指定的人物角色不存在") // 不存在的人物
	}

	// 虚拟角色存在才进行下一步
	friendModel, err := s.store.getOrCreate(ctx, userID, characterID)
	if err != nil {
		return nil, err
	}

	dto := s.convertToDTO(friendModel)
	return &dto, nil
}