package character

import "context"


func (cs *characterService) Exist(ctx context.Context, characterId uint) (bool, error){
	if characterId == 0 {
		return false, nil
	}

	return cs.characterStore.exist(ctx, characterId)
}