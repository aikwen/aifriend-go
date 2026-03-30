package meilisearch

import (
	"context"
	"log"
)

func (c *meilisearchClient[T]) Upsert(ctx context.Context, doc map[string]any) error {
	task, err := c.client.Index(c.index).UpdateDocumentsWithContext(ctx, []map[string]any{doc}, nil)
	if err != nil {
		log.Printf("[Meilisearch] 索引 %s 同步失败: %v\n", c.index, err)
		return err
	}

	log.Printf("[Meilisearch] 索引 %s 提交同步任务成功, TaskUID: %d\n", c.index, task.TaskUID)
	return nil
}


// UpsertInBatches 批量导入数据
func (c *meilisearchClient[T]) UpsertInBatches(ctx context.Context, docs []map[string]any, batchSize int) error {
	if len(docs) == 0 {
		return nil
	}

	tasks, err := c.client.Index(c.index).UpdateDocumentsInBatchesWithContext(ctx, docs, batchSize, nil)
	if err != nil {
		log.Printf("[Meilisearch] 索引 %s 批量同步失败: %v\n", c.index, err)
		return err
	}

	for _, task := range tasks {
		log.Printf("[Meilisearch] 索引 %s 提交批量同步任务成功, TaskUID: %d\n", c.index, task.TaskUID)
	}

	return nil
}