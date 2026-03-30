package meilisearch

import (
	"context"
	"testing"
)

func TestUpsert(t *testing.T) {
	c := NewClient[uint](&MeilisearchConfig{
		Host: "192.168.62.131",
		Port: 7700,
		APIKey: "123456",
	},"character", "id")
	err := c.SetupIndex()
	if err != nil {
		t.Log(err)
	}
	case1 := map[string]any{"id":1, "profile":"我是齐天大圣", "name":"孙悟空", "age":"18"}
	case2 := map[string]any{"id":2, "profile":"我是天鹏元帅", "name":"猪八戒"}

	err = c.Upsert(context.Background(), case1)
	if err != nil {
		t.Log(err)
	}

	err = c.Upsert(context.Background(), case2)
	if err != nil {
		t.Log(err)
	}
}