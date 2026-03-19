package rag

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/aikwen/aifriend-go/internal/chat/llm/rag/embeddings"
)

func TestIngestFile(t *testing.T) {
	ctx := context.Background()

	apiKey := os.Getenv("API_KEY")
	apiBase := os.Getenv("API_BASE")
	embeddingModel := "text-embedding-v4"
	qdrantAPIKey := os.Getenv("QDRANT_API_KEY")
	qdrantHost := "127.0.0.1"
	qdrantGRPCPortStr := "6334"
	qdrantCollection := "my_knowledge"

	if apiKey == "" || apiBase == "" || embeddingModel == "" {
		t.Fatal("API_KEY / API_BASE / EMBEDDING_MODEL 不能为空")
	}
	if qdrantHost == "" || qdrantGRPCPortStr == "" || qdrantCollection == "" {
		t.Fatal("QDRANT_HOST / QDRANT_GRPC_PORT / QDRANT_COLLECTION 不能为空")
	}

	qdrantGRPCPort, err := strconv.Atoi(qdrantGRPCPortStr)
	if err != nil {
		t.Fatalf("parse QDRANT_GRPC_PORT failed: %v", err)
	}

	embedder, err := embeddings.NewEinoSvc(ctx, embeddings.EinoConfig{
		APIKey:     apiKey,
		BaseURL:    apiBase,
		Model:      embeddingModel,
		Dimensions: 1024,
		BatchSize:  10,
	})
	if err != nil {
		t.Fatalf("new eino embedder failed: %v", err)
	}

	svc, err := NewService(embedder, RagConfig{
		APIKey:     qdrantAPIKey,
		Host:       qdrantHost,
		GRPCPort:   qdrantGRPCPort,
		Collection: qdrantCollection,
		TopK:       3,
		VectorSize: 1024,
	})
	if err != nil {
		t.Fatalf("new rag service failed: %v", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("get user home dir failed: %v", err)
	}


	filePath := filepath.Join(homeDir,"ragtest", "ali.txt")

	if err := svc.IngestFile(ctx, filePath); err != nil {
		t.Fatalf("ingest file failed: %v", err)
	}

	result, err := svc.Search(ctx, "阿里云百炼收费制度")
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if result == "" {
		t.Fatal("search result is empty after ingest")
	}

	t.Logf("search result:\n%s", result)
}