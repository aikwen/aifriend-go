package embeddings

import (
	"context"
	"fmt"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
)

func (o *openaiSvc) EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return [][]float32{}, nil
	}

	allEmbeddings := make([][]float32, 0, len(texts))

	for i := 0; i < len(texts); i += o.batchSize {
		end := i + o.batchSize
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

		params := openai.EmbeddingNewParams{
			Model: o.model,
			Input: openai.EmbeddingNewParamsInputUnion{
				OfArrayOfStrings: batch,
			},
		}

		if o.dimensions > 0 {
			params.Dimensions = param.NewOpt(int64(o.dimensions))
		}

		resp, err := o.client.Embeddings.New(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("create embeddings failed at batch %d: %w", i, err)
		}

		for _, item := range resp.Data {
			row := make([]float32, len(item.Embedding))
			for j, v := range item.Embedding {
				row[j] = float32(v)
			}
			allEmbeddings = append(allEmbeddings, row)
		}
	}

	return allEmbeddings, nil
}