package config

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)


var GlobalConfig = &Config{}


// Config 全局配置根结构
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	DB     DBConfig     `mapstructure:"db"`
	JWT    JWTConfig    `mapstructure:"jwt"`
	Prometheus PrometheusConfig `mapstructure:"prometheus"`
	LLM        LLMConfig        `mapstructure:"llm"`
	Qdrant QdrantConfig  `mapstructure:"qdrant"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
	AllowRegister bool  `mapstructure:"allow_register"`
	Enable        bool   `mapstructure:"enable"`
}

type DBConfig struct {
	DsnMysql     string `mapstructure:"mysql_dsn"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int `mapstructure:"conn_max_lifetime"`
	ConnectMaxRetries int `mapstructure:"connect_max_retries"`
    ConnectRetryInterval int `mapstructure:"connect_retry_interval"`
}

type JWTConfig struct {
	AccessSecret           string `mapstructure:"access_secret"`
	RefreshSecret          string `mapstructure:"refresh_secret"`
	RotateRefreshTokens    bool   `mapstructure:"rotate_refresh_tokens"`
	BlacklistAfterRotation bool   `mapstructure:"blacklist_after_rotation"`
}

type PrometheusConfig struct {
	Enable   bool   `mapstructure:"enable"`
	HttpAddr string `mapstructure:"http_addr"`
}

type LLMConfig struct {
    APIKey              string `mapstructure:"api_key"`
    APIBase             string `mapstructure:"api_base"`
	ModelName           string `mapstructure:"model_name"`
	EmbeddingModel      string `mapstructure:"embedding_model"`
	EmbeddingDimensions int `mapstructure:"embedding_dimensions"`
	EmbeddingBatchSize  int `mapstructure:"embedding_batchsize"`
}

type QdrantConfig struct {
	APIKey      string `mapstructure:"api_key"`
	Host        string `mapstructure:"host"`
	HTTPPort   int    `mapstructure:"http_port"`
	GRPCPort   int    `mapstructure:"grpc_port"`
	Collection  string `mapstructure:"collection"`
	TopK           int    `mapstructure:"top_k"`
}

// LoadConfig 读取配置并组装成 Struct
func LoadConfig(configPath string) *Config {
	if err := godotenv.Load(); err != nil {
        log.Println("未找到 .env 文件，将仅使用 config.yaml 和系统环境变量")
    }

	v := viper.New()

	// 设置默认值
	v.SetDefault("server.port", ":8000")
	v.SetDefault("server.mode", "dev")
	v.SetDefault("server.enable", true)

	v.SetDefault("prometheus.enable", true)
	v.SetDefault("prometheus.http_addr", "127.0.0.1:8001")

	// 读取系统环境变量
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取指定的配置文件
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigFile(filepath.Join("config", "config.yaml"))
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("警告: 未找到 config.yaml (%v)", err)
		} else {
			log.Fatalf("读取配置文件失败: %v", err)
		}
	}

	newConf := &Config{}

	if err := v.Unmarshal(newConf); err != nil {
		log.Fatalf("配置解析失败: %v", err)
	}

	GlobalConfig = newConf

	// 热更新
	v.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("配置文件已修改: %s", e.Name)
		updatedConf := &Config{}
		if err := v.Unmarshal(updatedConf); err != nil {
			log.Printf("热更新配置解析失败: %v", err)
		} else {
			GlobalConfig = updatedConf
			log.Printf("配置已自动更新: %+v", GlobalConfig.Server)
		}
	})
	v.WatchConfig()

	return GlobalConfig
}