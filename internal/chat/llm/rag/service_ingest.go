package rag

import (
	"context"
	"fmt"
	"os"
	"strings"
	"github.com/aikwen/aifriend-go/pkg/hash"

	qdrant "github.com/qdrant/go-client/qdrant"
)

const (
	defaultChunkSize    = 300
	defaultChunkOverlap = 30
)

func (r *ragSvc) IngestFile(ctx context.Context, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read file failed: %w", err)
	}

	text := strings.TrimSpace(string(content))
	if text == "" {
		return fmt.Errorf("file content is empty")
	}

	chunks := SplitText(text, defaultChunkSize, defaultChunkOverlap)
	if len(chunks) == 0 {
		return fmt.Errorf("split text got empty chunks")
	}

	vectors, err := r.embedder.EmbedDocuments(ctx, chunks)
	if err != nil {
		return fmt.Errorf("embed documents failed: %w", err)
	}

	if len(vectors) != len(chunks) {
		return fmt.Errorf("embedding count mismatch: chunks=%d vectors=%d", len(chunks), len(vectors))
	}

	points := make([]*qdrant.PointStruct, 0, len(chunks))
	for i, chunk := range chunks {
		idKey := fmt.Sprintf("%s#%d", filePath, i)
		pointID := hash.Hash64([]byte(idKey))
		payload := map[string]*qdrant.Value{
			"content": {
				Kind: &qdrant.Value_StringValue{
					StringValue: chunk,
				},
			},
			"source": {
				Kind: &qdrant.Value_StringValue{
					StringValue: filePath,
				},
			},
			"chunk_index": {
				Kind: &qdrant.Value_IntegerValue{
					IntegerValue: int64(i),
				},
			},
		}

		points = append(points, &qdrant.PointStruct{
			Id: &qdrant.PointId{
				PointIdOptions: &qdrant.PointId_Num{
					Num: pointID,
				},
			},
			Vectors: &qdrant.Vectors{
				VectorsOptions: &qdrant.Vectors_Vector{
					Vector: &qdrant.Vector{
						Data: vectors[i],
					},
				},
			},
			Payload: payload,
		})
	}

	_, err = r.client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: r.collection,
		Points:         points,
	})
	if err != nil {
		return fmt.Errorf("qdrant upsert failed: %w", err)
	}

	return nil
}