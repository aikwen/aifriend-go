package memory

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/chat/llm/graph"
	"github.com/aikwen/aifriend-go/internal/store"
	"github.com/aikwen/aifriend-go/internal/store/models"
)


type Service interface {
	Update(ctx context.Context, userID uint, friend *models.Friend) error
}


type memorySvc struct {
	graph    *graph.Graph
	database *store.Database
}

func NewService(graph *graph.Graph, database *store.Database) Service {
	return &memorySvc{
		graph: graph,
		database: database,
	}
}