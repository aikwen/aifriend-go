package main

import (
	"log"

	"github.com/aikwen/aifriend-go/config"

	"github.com/aikwen/aifriend-go/internal/api/router"
	"github.com/aikwen/aifriend-go/internal/api/handler"
	"github.com/aikwen/aifriend-go/internal/models"
	"github.com/aikwen/aifriend-go/internal/auth"
	"github.com/aikwen/aifriend-go/internal/user"
	"github.com/aikwen/aifriend-go/internal/character"
	"github.com/aikwen/aifriend-go/pkg/db"
	"github.com/aikwen/aifriend-go/pkg/storage"
)

func main(){
	// 加载环境变量
	cfg := config.LoadConfig()
	// 数据库
	gormDB := db.InitMySQL(cfg.DB, cfg.Server.Mode)

	log.Println("正在进行数据库迁移...")
	if err := gormDB.AutoMigrate(&models.User{}, &models.Character{}); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	log.Println("数据库迁移结束...")

	// 依赖注入
	fileStorage := storage.NewLocalStorage("media")
	userSvc := user.NewUserService(gormDB, fileStorage)
	charSvc := character.NewCharacterService(gormDB, fileStorage)
	authSvc := auth.NewAuthService(userSvc, &cfg.JWT)
	h := handler.NewHandler(authSvc, charSvc, userSvc, fileStorage)
	r := router.SetupRouter(h, cfg.JWT.AccessSecret, cfg.Server.Mode)

	// 启动
	log.Printf("🚀 服务启动成功！运行环境: %s, 监听端口: %s", cfg.Server.Mode, cfg.Server.Port)
	if err := r.Run(cfg.Server.Port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}