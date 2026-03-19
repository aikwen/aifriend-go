package embeddings

import (
	"context"
	"os"
	"testing"
)

func TestOpenAIEmbedQuery(t *testing.T) {
	ctx := context.Background()

	apiKey := os.Getenv("API_KEY")
	apiBase := os.Getenv("API_BASE")
	embeddingModel := os.Getenv("EMBEDDING_MODEL")

	if apiKey == "" || apiBase == "" || embeddingModel == "" {
		t.Fatal("API_KEY / API_BASE / EMBEDDING_MODEL 不能为空")
	}

	svc, err := NewOpenAISvc(OpenAIConfig{
		APIKey:     apiKey,
		BaseURL:    apiBase,
		Model:      embeddingModel,
		Dimensions: 1024,
		BatchSize:  10,
	})
	if err != nil {
		t.Fatalf("new openai embedding svc failed: %v", err)
	}

	vec, err := svc.EmbedQuery(ctx, "阿里云百炼平台是什么？")
	if err != nil {
		t.Fatalf("embed query failed: %v", err)
	}

	if len(vec) == 0 {
		t.Fatal("embedding result is empty")
	}

	t.Logf("embedding length=%d", len(vec))

	preview := 5
	if len(vec) < preview {
		preview = len(vec)
	}
	t.Logf("embedding preview=%v", vec[:preview])
}