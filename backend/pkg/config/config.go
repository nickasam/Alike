// Package config 从环境变量加载应用配置，支持 .env 文件，所有字段带默认值。
package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config 聚合应用运行所需的全部配置。
type Config struct {
	Env        string // 运行环境：development/production
	ServerPort string // HTTP 服务端口

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	RedisAddr     string
	RedisPassword string
	RedisDB       int

	JWTSecret        string
	JWTAccessExpire  time.Duration // access token 有效期
	JWTRefreshExpire time.Duration // refresh token 有效期
}

// Load 优先加载 .env（若存在）到环境变量，再从环境变量读取配置。
func Load() *Config {
	loadDotEnv(".env")

	return &Config{
		Env:        getEnv("APP_ENV", "development"),
		ServerPort: getEnv("SERVER_PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "alike"),
		DBPassword: getEnv("DB_PASSWORD", "alike"),
		DBName:     getEnv("DB_NAME", "alike"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),

		JWTSecret:        getEnv("JWT_SECRET", "alike-dev-secret-change-me"),
		JWTAccessExpire:  time.Duration(getEnvInt("JWT_ACCESS_EXPIRE_MIN", 120)) * time.Minute,
		JWTRefreshExpire: time.Duration(getEnvInt("JWT_REFRESH_EXPIRE_HOUR", 168)) * time.Hour,
	}
}

// DSN 返回 PostgreSQL 连接字符串（pgx stdlib / database/sql 使用）。
func (c *Config) DSN() string {
	return "host=" + c.DBHost +
		" port=" + c.DBPort +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" sslmode=" + c.DBSSLMode
}

// loadDotEnv 解析简单的 KEY=VALUE 文件，已存在的环境变量不覆盖。
func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return // .env 不存在时静默跳过
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.Trim(strings.TrimSpace(val), `"'`)
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, val)
		}
	}
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v, ok := os.LookupEnv(key); ok {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
