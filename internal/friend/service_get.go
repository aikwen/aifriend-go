package friend


func (s *friendService) GetList(userID uint, offset int, limit int) ([]FriendDTO, error) {
	friendModels, err := s.store.getList(userID, offset, limit)
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