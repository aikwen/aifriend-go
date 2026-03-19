package rag

import (
	"context"
	"fmt"
	"strings"

	qdrant "github.com/qdrant/go-client/qdrant"
)

func (r *ragSvc) Search(ctx context.Context, query string) (string, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return "", fmt.Errorf("query is empty")
	}

	vector, err := r.embedder.EmbedQuery(ctx, query)
	if err != nil {
		return "", fmt.Errorf("embed query failed: %w", err)
	}

	limit := uint64(r.topK)

	resp, err := r.client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: r.collection,
		Query:          qdrant.NewQuery(vector...),
		Limit:          &limit,
		WithPayload:    qdrant.NewWithPayload(true),
	})
	if err != nil {
		return "", fmt.Errorf("qdrant query failed: %w", err)
	}

	if len(resp) == 0 {
		return "", nil
	}

	var builder strings.Builder
	builder.WriteString("从知识库找到以下相关信息：\n\n")

	count := 0
	for _, point := range resp {
		if point == nil || point.Payload == nil {
			continue
		}

		val, ok := point.Payload["content"]
		if !ok || val == nil {
			continue
		}

		content := val.GetStringValue()
		content = strings.TrimSpace(content)
		if content == "" {
			continue
		}

		count++
		fmt.Fprintf(&builder, "内容片段%d：\n", count)
		fmt.Fprint(&builder, content)
		fmt.Fprint(&builder, "\n\n")
	}

	if count == 0 {
		return "", nil
	}

	return strings.TrimSpace(builder.String()), nil
}