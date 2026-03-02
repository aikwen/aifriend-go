package friend

import "context"


func (s *friendService) GetList(ctx context.Context, userID uint, offset int, limit int) ([]FriendDTO, error) {
	friendModels, err := s.store.getList(ctx, userID, offset, limit)
	if err != nil {
		return nil, err
	}

	var dtoList []FriendDTO
	for _, f := range friendModels {
		dtoList = append(dtoList, s.convertToDTO(&f))
	}

	if dtoList == nil {
		dtoList = make([]FriendDTO, 0)
	}

	return dtoList, nil
}