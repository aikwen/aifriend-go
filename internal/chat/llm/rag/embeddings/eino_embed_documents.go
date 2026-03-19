package embeddings

import (
	"context"
	"fmt"
	"strings"
)

// EmbedDocuments 将 texts 文本块转换成向量
func (e *einoSvc) EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return [][]float32{}, nil
	}

	allEmbeddings := make([][]float32, 0, len(texts))

	for i := 0; i < len(texts); i += e.batchSize {
		end := i + e.batchSize
		if end > len(texts) {
			end = len(texts)
		}

		rawBatch := texts[i:end]
		batch := make([]string, 0, len(rawBatch))
		for _, text := range rawBatch {
			text = strings.TrimSpace(text)
			if text != "" {
				batch = append(batch, text)
			}
		}

		if len(batch) == 0 {
			continue
		}

		vectors, err := e.embedder.EmbedStrings(ctx, batch)
		if err != nil {
			return nil, fmt.Errorf("eino embed documents failed at batch %d: %w", i, err)
		}

		for _, vec := range vectors {
			row := make([]float32, len(vec))
			for j, v := range vec {
				row[j] = float32(v)
			}
			allEmbeddings = append(allEmbeddings, row)
		}
	}

	return allEmbeddings, nil
}