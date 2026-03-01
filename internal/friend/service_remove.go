package friend


func (s *friendService) Remove(userID uint, friendID uint) error {
	return s.store.remove(userID, friendID)
}