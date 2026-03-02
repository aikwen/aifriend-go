package friend

import "context"


func (s *friendService) Remove(ctx context.Context, userID uint, friendID uint) error {
	return s.store.remove(ctx, userID, friendID)
}