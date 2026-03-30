package main

import (
	"context"
	"flag"
	"log"

	"github.com/aikwen/aifriend-go/config"
	"github.com/aikwen/aifriend-go/internal/store/models"
	"github.com/aikwen/aifriend-go/pkg/db"
	"github.com/aikwen/aifriend-go/pkg/meilisearch"
	"gorm.io/gorm"
)

func main() {
	configFile := flag.String("config-file", "config/config.yaml", "配置文件路径")
	clean := flag.Bool("clean", false, "同步前是否清空 Meilisearch 中的所有旧数据")
	flag.Parse()


	cfg := config.LoadConfig(*configFile)

	gormDB := db.InitMySQL(cfg.DB, cfg.Server.Mode)

	msCfg := &meilisearch.MeilisearchConfig{
		Host:   cfg.MeiliSearch.Host,
		Port:   cfg.MeiliSearch.Port,
		APIKey: cfg.MeiliSearch.APIKey,
	}

	// 初始化客户端，主键类型为 uint
	msClient := meilisearch.NewClient[uint](msCfg, "characters", "ID")

	if err := msClient.SetupIndex(); err != nil {
		log.Fatalf("设置索引失败: %v", err)
	}

	if *clean {
		log.Println("检测到 -clean 参数，正在清空 Meilisearch 旧数据...")
		if err := msClient.DeleteAllDocuments(); err != nil {
			log.Fatalf("清空旧数据失败: %v", err)
		}
	}

	ctx := context.Background()
	batchSize := 100 // 每次从数据库提取并在 Meilisearch 中更新的条数

	// 用于接收当前批次数据的切片
	var characters []models.Character

	log.Printf("开始全量同步，每次处理 %d 条数据...", batchSize)

	result := gormDB.Select("id", "name", "profile").
		Where("deleted_at IS NULL").
		FindInBatches(&characters, batchSize, func(tx *gorm.DB, batch int) error {

			// 预先分配好容量，提升性能
			docs := make([]map[string]any, 0, len(characters))

			// 2. 手动构建文档映射，完全抛弃不需要的字段
			for _, char := range characters {
				docs = append(docs, map[string]any{
					"ID":      char.ID, // 注意：这里的键名需与 msClient 初始化时传入的主键名一致
					"name":    char.Name,
					"profile": char.Profile,
				})
			}

			// 3. 提交给 Meilisearch
			if err := msClient.UpsertInBatches(ctx, docs, batchSize); err != nil {
				log.Printf("第 %d 批次同步失败: %v", batch, err)
				return err
			}

			log.Printf("第 %d 批次同步成功", batch)
			return nil
		})

	if result.Error != nil {
		log.Fatalf("全量同步过程中发生错误: %v", result.Error)
	}

	log.Printf("全量同步完成！共处理数据行数: %d", result.RowsAffected)

}