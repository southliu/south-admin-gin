package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

// Config 应用配置
type Config struct {
	MySQL MySQLConfig `yaml:"mysql"`
	JWT   JWTConfig   `yaml:"jwt"`
	BFB   BFBConfig   `yaml:"bfb"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret     string `yaml:"secret"`
	ExpireHour int    `yaml:"expire-hour"`
}

// MySQLConfig MySQL 数据库配置
type MySQLConfig struct {
	Prefix            string `yaml:"prefix"`
	Port              string `yaml:"port"`
	Config            string `yaml:"config"`
	DBName            string `yaml:"db-name"`
	Username          string `yaml:"username"`
	Password          string `yaml:"password"`
	Path              string `yaml:"path"`
	Engine            string `yaml:"engine"`
	LogMode           string `yaml:"log-mode"`
	MaxIdleConns      int    `yaml:"max-idle-conns"`
	MaxOpenConns      int    `yaml:"max-open-conns"`
	ConnMaxLifetime   int    `yaml:"conn-max-lifetime"`
	Singular          bool   `yaml:"singular"`
	LogZap            bool   `yaml:"log-zap"`
}

// BFBConfig BFB 相关配置
type BFBConfig struct {
	DingTalk DingTalkConfig `yaml:"dingtalk"`
	Sync     SyncConfig     `yaml:"sync"`
}

// DingTalkConfig 钉钉配置
type DingTalkConfig struct {
	AppKey           string `yaml:"app-key"`
	AppSecret        string `yaml:"app-secret"`
	RedirectURI      string `yaml:"redirect-uri"`
	FrontendCallback string `yaml:"frontend-callback"`
	ProcessCode      string `yaml:"process-code"`
}

// SyncConfig 数据同步配置
type SyncConfig struct {
	APIKey string `yaml:"api-key"`
}

var globalConfig *Config

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	globalConfig = &cfg
	return &cfg, nil
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	return globalConfig
}

// GetDSN 获取 MySQL 连接字符串
func (c *MySQLConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		c.Username,
		c.Password,
		c.Path,
		c.Port,
		c.DBName,
		c.Config,
	)
}
