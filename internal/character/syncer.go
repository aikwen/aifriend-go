package character

import (
	"context"
	"log"

	"github.com/aikwen/aifriend-go/internal/store/models"
	"github.com/aikwen/aifriend-go/pkg/meilisearch"
	"github.com/aikwen/aifriend-go/pkg/task"
	"gorm.io/gorm"
)

// MeiliSyncer 异步同步器
type MeiliSyncer struct {
	msClient meilisearch.Client[uint]
	db       *gorm.DB
	channels [10]chan uint
}

func NewMeiliSyncer(client meilisearch.Client[uint], db *gorm.DB) *MeiliSyncer {
	syncer := &MeiliSyncer{
		msClient: client,
		db:       db,
	}

	for i := range 10 {
		syncer.channels[i] = make(chan uint, 1024)
		go syncer.worker(i)
	}

	return syncer
}

func (s *MeiliSyncer) worker(workerID int) {
	ctx := context.Background()

	for id := range s.channels[workerID] {
		var char models.Character

		err := s.db.Unscoped().
			Select("id", "name", "profile", "deleted_at").
			First(&char, id).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// 如果数据真的在物理层面不存在了，可以选择直接在 Meili 中物理删除
				_ = s.msClient.Upsert(ctx, map[string]any{"ID": id, "deleted_at": 9999999999})
			} else {
				log.Printf("[MeiliSyncer] Worker-%d 查询角色(ID:%d)失败: %v\n", workerID, id, err)
			}
			continue
		}

		// 同步内容
		doc := map[string]any{
			"ID":      char.ID,
			"name":    char.Name,
			"profile": char.Profile,
		}

		if char.DeletedAt.Valid {
			doc["deleted_at"] = char.DeletedAt.Time.Unix()
		} else {
			// 正常更新
			doc["deleted_at"] = nil
		}

		// 推送 Meilisearch
		if err := s.msClient.Upsert(ctx, doc); err != nil {
			log.Printf("[MeiliSyncer] Worker-%d 同步角色(ID:%d)失败: %v\n", workerID, char.ID, err)
		}
	}
}

// Enqueue 用来更新character 为 id 的信息
func (s *MeiliSyncer) Enqueue(id uint) {
	if id == 0 {
		return
	}

	idx := id % 10
	task.Go(func() {
		s.channels[idx] <- id
	})
}