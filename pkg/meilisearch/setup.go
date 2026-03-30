package meilisearch

import (
	"log"
)


// SetupIndex 配置索引的过滤和排序规则（在系统启动时或迁移脚本中调用一次即可）
func (c *meilisearchClient[T]) SetupIndex() error {

	task1, err := c.client.Index(c.index).UpdateFilterableAttributes(&[]any{"deleted_at"})
	if err != nil {
		log.Printf("设置可过滤字段失败: %v", err)
		return err
	}


	task2, err := c.client.Index(c.index).UpdateSortableAttributes(&[]string{c.primaryKey})
	if err != nil {
		log.Printf("设置可排序字段失败: %v", err)
		return err
	}

	log.Printf("[Meilisearch] 索引 %s 规则配置已提交 (Task: %d, %d)\n", c.index, task1.TaskUID, task2.TaskUID)
	return nil
}