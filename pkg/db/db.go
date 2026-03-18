package db

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/aikwen/aifriend-go/config"
	"github.com/aikwen/aifriend-go/pkg/monitor"
)

// InitMySQL 初始化 MySQL 连接
// dsn "root:123456@tcp(127.0.0.1:3306)/aifriends_db?charset=utf8mb4&parseTime=True&loc=Local"
// env "prod", "dev"
func InitMySQL(cfg config.DBConfig, appEnv string) *gorm.DB {
	// 配置日志
	logLevel := logger.Info
	if appEnv == "prod" {
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

	var db *gorm.DB
	var err error
	maxRetries := cfg.ConnectMaxRetries
	retryInterval := cfg.ConnectRetryInterval
	if maxRetries <= 0 {
		maxRetries = 1
	}

	if retryInterval <= 0{
		retryInterval = 3
	}

	for i := range maxRetries {
		// 打开连接
		db, err = gorm.Open(mysql.Open(cfg.DsnMysql), &gorm.Config{
			Logger: newLogger,
		})

		if err == nil {
			sqlDB, dbErr := db.DB()
			if dbErr == nil {
				pingErr := sqlDB.Ping()
				if pingErr == nil {
					log.Println("连接 MySQL 成功")
					break
				}
				err = pingErr
			} else {
				err = dbErr
			}
		}

		log.Printf("连接 MySQL 失败（%d/%d）: %v", i+1, maxRetries, err)

		if i == maxRetries-1 {
			log.Fatalf("连接 MySQL 最终失败: %v", err)
		}

		time.Sleep(time.Second * time.Duration(retryInterval))
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("获取底层 sql.DB 失败: %v", err)
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	// SetMaxOpenConns 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	// SetConnMaxLifetime 设置连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(cfg.ConnMaxLifetime))

	// 添加插件
	if err := db.Use(&monitor.GormMetrics{}); err != nil {
		log.Fatalf("注册 GORM 监控插件失败: %v", err)
	}
	log.Println("GORM Prometheus 监控插件挂载成功")
	return db
}