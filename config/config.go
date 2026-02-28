package config


import (
	"log"
	"strings"

	"github.com/spf13/viper"
)


// Config 全局配置根结构
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	DB     DBConfig     `mapstructure:"db"`
	JWT    JWTConfig    `mapstructure:"jwt"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
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

// LoadConfig 读取配置并组装成 Struct
func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// 设置默认值
	viper.SetDefault("server.port", ":8000")
	viper.SetDefault("server.mode", "dev")

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

	return &config
}