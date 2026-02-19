package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/aikwen/aifriend-go/internal/handler"
	"github.com/aikwen/aifriend-go/internal/models"
	"github.com/aikwen/aifriend-go/internal/router"
	"github.com/aikwen/aifriend-go/internal/service"
	"github.com/aikwen/aifriend-go/internal/store"
	"github.com/aikwen/aifriend-go/pkg/db"
)

func main(){
	//加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("未找到 .env 文件，将尝试使用系统环境变量")
	}
	//加载数据库环境变量
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		log.Fatal("环境变量 MYSQL_DSN 未设置，请检查 .env 文件")
	}
	// JWT 环境变量
	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if accessSecret == "" || refreshSecret == "" {
		log.Fatal("JWT 密钥 (ACCESS_SECRET/REFRESH_SECRET) 未设置，请检查 .env")
	}

	rotateStr := os.Getenv("JWT_ROTATE_REFRESH_TOKENS")
	rotate, err := strconv.ParseBool(rotateStr)
	if err != nil {
		rotate = false
		log.Println("提示: JWT_ROTATE_REFRESH_TOKENS 未设置或无效，默认关闭 Token 轮换")
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "dev" //默认开发环境
	}
	//数据库
	gormDB := db.InitMySQL(dsn, appEnv)

	log.Println("正在进行数据库迁移...")
	if err := gormDB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	log.Println("数据库迁移结束...")
	// 依赖注入
	userStore := store.NewUserStore(gormDB)
	authSvc := service.NewAuthService(userStore, accessSecret, refreshSecret, rotate)
	userSvc := service.NewUserService(userStore)

	h := handler.NewHandler(authSvc, userSvc)
	r := router.SetupRouter(h, accessSecret, appEnv)

	//启动
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = ":8000"
	}
	log.Printf("🚀 服务启动成功！运行环境: %s, 监听端口: %s", appEnv, serverPort)
	if err := r.Run(serverPort); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}