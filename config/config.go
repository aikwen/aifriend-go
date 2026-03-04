package config

import (
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config 全局配置根结构
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	DB     DBConfig     `mapstructure:"db"`
	JWT    JWTConfig    `mapstructure:"jwt"`
	Prometheus PrometheusConfig `mapstructure:"prometheus"`
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

// LoadConfig 读取配置并组装成 Struct
func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// 设置默认值
	viper.SetDefault("server.port", ":8000")
	viper.SetDefault("server.mode", "dev")

	viper.SetDefault("prometheus.enable", true)
	viper.SetDefault("prometheus.http_addr", "127.0.0.1:8001")

	viper.SetDefault("server.enable", true)
	
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("警告: 未找到 config.yaml 配置文件，将完全依赖环境变量")
		} else {
			log.Fatalf("读取配置文件失败: %v", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("配置解析失败: %v", err)
	}

	// 热更新
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("配置文件已修改: %s", e.Name)
		if err := viper.Unmarshal(&config); err != nil {
			log.Printf("热更新配置解析失败: %v", err)
		} else {
			log.Printf("配置已自动更新: %+v", config.Server)
		}
	})
	viper.WatchConfig()

	return &config
}