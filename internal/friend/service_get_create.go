package friend


func (s *friendService) GetOrCreate(userID uint, characterID uint) (*FriendDTO, error) {
	friendModel, err := s.store.getOrCreate(userID, characterID)
	if err != nil {
		return nil, err
	}

	dto := s.convertToDTO(friendModel)
	return &dto, nil
}