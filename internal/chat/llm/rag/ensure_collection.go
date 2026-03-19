package rag

import (
	"context"
	"fmt"
	"strings"

	qdrant "github.com/qdrant/go-client/qdrant"
)

func ensureCollection(ctx context.Context, client *qdrant.Client, collection string, vectorSize int) error {
	if client == nil {
		return fmt.Errorf("qdrant client is nil")
	}
	if strings.TrimSpace(collection) == "" {
		return fmt.Errorf("qdrant collection is empty")
	}
	if vectorSize <= 0 {
		return fmt.Errorf("qdrant vector size must be greater than 0")
	}

	exists, err := client.CollectionExists(ctx, collection)
	if err != nil {
		return fmt.Errorf("check qdrant collection exists failed: %w", err)
	}
	if exists {
		return nil
	}

	err = client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: collection,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     uint64(vectorSize),
			Distance: qdrant.Distance_Cosine,
		}),
	})
	if err != nil {
		return fmt.Errorf("create qdrant collection failed: %w", err)
	}

	return nil
}