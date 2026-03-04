package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port int
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type JWTConfig struct {
	Secret            string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

// Load 从环境变量和配置文件加载配置
func Load() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AutomaticEnv()

	// 设置默认值
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		// 如果配置文件不存在，尝试使用环境变量
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
	}

	// 构建配置
	cfg := &Config{
		Server: ServerConfig{
			Port: viper.GetInt("SERVER_PORT"),
			Mode: viper.GetString("SERVER_MODE"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetInt("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSL_MODE"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetInt("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		JWT: JWTConfig{
			Secret:            viper.GetString("JWT_SECRET"),
			AccessTokenExpiry:  viper.GetDuration("JWT_ACCESS_TOKEN_EXPIRY"),
			RefreshTokenExpiry: viper.GetDuration("JWT_REFRESH_TOKEN_EXPIRY"),
		},
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return cfg, nil
}

// setDefaults 设置默认配置值
func setDefaults() {
	viper.SetDefault("SERVER_PORT", 8080)
	viper.SetDefault("SERVER_MODE", "development")

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_NAME", "alike_db")
	viper.SetDefault("DB_USER", "alike_user")
	viper.SetDefault("DB_PASSWORD", "")
	viper.SetDefault("DB_SSL_MODE", "disable")

	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", 6379)
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)

	viper.SetDefault("JWT_SECRET", "change-this-in-production")
	viper.SetDefault("JWT_ACCESS_TOKEN_EXPIRY", 15*time.Minute)
	viper.SetDefault("JWT_REFRESH_TOKEN_EXPIRY", 168*time.Hour)
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("无效的服务器端口: %d", c.Server.Port)
	}

	if c.Database.Host == "" {
		return fmt.Errorf("数据库主机不能为空")
	}

	if c.Database.DBName == "" {
		return fmt.Errorf("数据库名称不能为空")
	}

	if c.JWT.Secret == "" || c.JWT.Secret == "change-this-in-production" {
		if os.Getenv("SERVER_MODE") == "production" {
			return fmt.Errorf("JWT密钥必须在生产环境中设置")
		}
	}

	return nil
}

// GetDSN 生成数据库连接字符串
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai",
		c.Database.Host,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.Port,
		c.Database.SSLMode,
	)
}

// GetRedisAddr 生成Redis连接地址
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}
