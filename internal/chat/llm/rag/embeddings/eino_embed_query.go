package embeddings

import (
	"context"
	"fmt"
)

func (e *einoSvc) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	embeddings, err := e.EmbedDocuments(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	if len(embeddings) == 0 {
		return nil, fmt.Errorf("query embedding is empty")
	}

	return embeddings[0], nil
}