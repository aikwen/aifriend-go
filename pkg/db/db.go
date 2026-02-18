package db

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitMySQL 初始化 MySQL 连接
// dsn "root:123456@tcp(127.0.0.1:3306)/aifriends_db?charset=utf8mb4&parseTime=True&loc=Local"
// env "prod", "dev"
func InitMySQL(dsn string, env string) *gorm.DB {
	// 配置日志
	logLevel := logger.Info
	if env == "prod" {
		logLevel = logger.Warn
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// 打开连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("连接 MySQL 失败: %v", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取底层 sql.DB 失败: %v", err)
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}