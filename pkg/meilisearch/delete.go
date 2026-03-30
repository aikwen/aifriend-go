package meilisearch

import "log"


func (c *meilisearchClient[T]) DeleteAllDocuments() error {
	task, err := c.client.Index(c.index).DeleteAllDocuments(nil)
	if err != nil {
		log.Printf("[Meilisearch] 索引 %s 清空失败: %v\n", c.index, err)
		return err
	}

	log.Printf("[Meilisearch] 索引 %s 提交清空任务成功, TaskUID: %d\n", c.index, task.TaskUID)
	return nil
}