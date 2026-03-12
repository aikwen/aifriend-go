package main

import (
	"log"

	"github.com/aikwen/aifriend-go/config"
	"github.com/aikwen/aifriend-go/internal/api/handler"
	"github.com/aikwen/aifriend-go/internal/api/router"
	"github.com/aikwen/aifriend-go/internal/auth"
	"github.com/aikwen/aifriend-go/internal/character"
	"github.com/aikwen/aifriend-go/internal/chat"
	"github.com/aikwen/aifriend-go/internal/friend"
	"github.com/aikwen/aifriend-go/internal/store"
	"github.com/aikwen/aifriend-go/internal/store/models"
	"github.com/aikwen/aifriend-go/internal/user"
	"github.com/aikwen/aifriend-go/pkg/db"
	"github.com/aikwen/aifriend-go/pkg/monitor"
	"github.com/aikwen/aifriend-go/pkg/storage"
)

func main() {
	monitor.Init()
	// 加载环境变量
	cfg := config.LoadConfig()
	// 数据库
	gormDB := db.InitMySQL(cfg.DB, cfg.Server.Mode)

	log.Println("正在进行数据库迁移...")
	if err := gormDB.AutoMigrate(&models.User{},
		&models.Character{},
		&models.Friend{},
		&models.Message{},
		//&models.SystemPrompt{},
		); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	log.Println("数据库迁移结束...")

	// 依赖注入
	fileStorage := storage.NewLocalStorage("media")
	database := store.NewDatabase(gormDB)
	userSvc := user.NewUserService(database, fileStorage)
	charSvc := character.NewCharacterService(database, fileStorage)
	authSvc := auth.NewAuthService(database)
	friendSvc := friend.NewFriendService(database)
	chatSvc, err := chat.NewChatService(database)
	if err != nil {
		log.Fatal("chatSvc init error", err)
	}
	h := handler.NewHandler(authSvc, charSvc, userSvc, friendSvc,chatSvc, fileStorage)
	// 初始化路由
	r := router.SetupRouter(h)
	// 启动 prometheus
	if cfg.Prometheus.Enable {
		go monitor.StartMetricsServer(cfg.Prometheus.HttpAddr)
	}

	// 启动
	log.Printf("🚀 服务启动成功！运行环境: %s, 监听端口: %s", cfg.Server.Mode, cfg.Server.Port)
	if err := r.Run(cfg.Server.Port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
