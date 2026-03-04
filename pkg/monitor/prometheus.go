package monitor

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)


var HttpRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "aifriend_http_requests_total",
		Help: "HTTP请求总数",
	},
	[]string{"method", "endpoint", "code"}, // 核心标签：请求方法、接口、状态码
)

var HttpRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "aifriend_http_request_duration_seconds",
		Help:    "HTTP请求耗时（单位：秒）",
		Buckets: []float64{.01, .05, .1, .25, .5, 1, 2.5, 5},
	},
	[]string{"method", "endpoint"},
)

var SqlExecDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "aifriend_sql_exec_duration_seconds",
		Help:    "数据库SQL执行耗时（单位：秒）",
		Buckets: []float64{.001, .005, .01, .05, .1, .5, 1},
	},
	[]string{"sql_type", "table", "success"},
)


func Init() {
	// 注册指标
	prometheus.MustRegister(
		HttpRequestsTotal,
		HttpRequestDuration,
		SqlExecDuration,
	)
}

func StartMetricsServer(addr string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[Panic] Prometheus 监控协程崩溃: %v\n堆栈信息: %s", r, debug.Stack())
		}
	}()

	mux := http.NewServeMux()
	// 挂载 Prometheus 指标接口
	mux.Handle("/metrics", promhttp.Handler())

	log.Printf("监控服务已启动，监听地址: http://%s/metrics (仅限本地访问)", addr)
	if err := http.ListenAndServe(addr, mux); err != nil && err != http.ErrServerClosed {
		log.Printf("[Error] Prometheus 监控服务器非正常退出: %v", err)
	}
}