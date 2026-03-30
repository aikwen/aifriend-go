package meilisearch

import (
	"context"
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {
	c := NewClient[uint](&MeilisearchConfig{
		Host: "192.168.62.131", Port: 7700, APIKey: "123456",
	}, "character", "id")

	cases := []string {
		"孙悟空",
		"齐天",
		"齐天大圣",
		"天鹏",
		"",
	}

	for i, s := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			res, err := c.Search(context.Background(), s, 10, 0)
			if err != nil {
				t.Logf("【%v】搜索报错:%v",i ,err)
			}
			t.Logf("【%v】搜索结果: %v", i, res)
		})
	}
}