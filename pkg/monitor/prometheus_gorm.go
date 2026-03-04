package monitor

import (
	"time"

	"gorm.io/gorm"
)


const gormInstanceKey = "startTime"


type GormMetrics struct{}


func (m *GormMetrics) Name() string {
	return "prometheus_metrics"
}

// Initialize 实现 GORM 插件接口，初始化并注册回调钩子
func (m *GormMetrics) Initialize(db *gorm.DB) error {
	// 1注册 Before 钩子：SQL 执行前记录开始时间
	_ = db.Callback().Create().Before("gorm:create").Register("metrics:before_create", before)
	_ = db.Callback().Query().Before("gorm:query").Register("metrics:before_query", before)
	_ = db.Callback().Update().Before("gorm:update").Register("metrics:before_update", before)
	_ = db.Callback().Delete().Before("gorm:delete").Register("metrics:before_delete", before)

	//  注册 After 钩子：SQL 执行后计算耗时并上报指标
	_ = db.Callback().Create().After("gorm:create").Register("metrics:after_create", after("insert"))
	_ = db.Callback().Query().After("gorm:query").Register("metrics:after_query", after("select"))
	_ = db.Callback().Update().After("gorm:update").Register("metrics:after_update", after("update"))
	_ = db.Callback().Delete().After("gorm:delete").Register("metrics:after_delete", after("delete"))

	return nil
}

// before 钩子将 SQL 执行开始时间存入 GORM 实例
func before(db *gorm.DB) {
	db.InstanceSet(gormInstanceKey, time.Now()) // 存入 time.Time 类型的当前时间
}

// after 钩子取出开始时间，计算耗时并上报 Prometheus 指标
func after(sqlType string) func(*gorm.DB) {
	return func(db *gorm.DB) {

		val, ok := db.InstanceGet(gormInstanceKey)
		if !ok {
			return
		}
		startTime, ok := val.(time.Time)
		if !ok {
			return
		}

		// 计算 SQL 真实执行耗时
		duration := time.Since(startTime).Seconds()

		// 判断 SQL 执行是否成功
		success := "true"
		if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
			success = "false"
		}

		tableName := db.Statement.Table
		if tableName == "" && db.Statement.Schema != nil {
			tableName = db.Statement.Schema.Table
		}
		if tableName == "" {
			tableName = "unknown_table"
		}

		SqlExecDuration.WithLabelValues(sqlType, tableName, success).Observe(duration)
	}
}