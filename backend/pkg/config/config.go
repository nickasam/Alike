// Package config 从环境变量加载应用配置，支持 .env 文件，所有字段带默认值。
package config

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

// defaultJWTSecret 是开发环境的占位密钥；生产环境必须显式覆盖，否则 Validate 报错。
const defaultJWTSecret = "alike-dev-secret-change-me"

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

	MinIOEndpoint  string // MinIO/S3 访问地址（host:port）
	MinIOAccessKey string // access key（MINIO_ROOT_USER）
	MinIOSecretKey string // secret key（MINIO_ROOT_PASSWORD）
	MinIOBucket    string // 默认 bucket
	MinIOUseSSL    bool   // 是否走 HTTPS
	MinIOPublicURL string // 对外可访问的基地址（用于拼接返回 URL），缺省用 endpoint

	UploadMaxImageBytes int64 // 图片上传上限（字节）
	UploadMaxDocBytes   int64 // 文档上传上限（字节）

	// CORSAllowedOrigins 是允许跨域的 Origin 白名单（精确匹配）。
	// 空表示：开发环境放开（回显请求 Origin），生产环境拒绝所有跨域。
	CORSAllowedOrigins []string
}

// Load 优先加载 .env（若存在）到环境变量，再从环境变量读取配置。
// 变量名与 docker-compose / .env 契约保持一致（POSTGRES_*/REDIS_HOST+PORT/BACKEND_PORT/JWT_*_TTL），
// 并兼容旧命名（DB_*/REDIS_ADDR/SERVER_PORT）作为回退。
func Load() *Config {
	loadDotEnv(".env")

	return &Config{
		Env:        getEnv("APP_ENV", "development"),
		ServerPort: firstEnv([]string{"BACKEND_PORT", "SERVER_PORT"}, "8080"),

		DBHost:     firstEnv([]string{"POSTGRES_HOST", "DB_HOST"}, "localhost"),
		DBPort:     firstEnv([]string{"POSTGRES_PORT", "DB_PORT"}, "5432"),
		DBUser:     firstEnv([]string{"POSTGRES_USER", "DB_USER"}, "alike"),
		DBPassword: firstEnv([]string{"POSTGRES_PASSWORD", "DB_PASSWORD"}, "alike"),
		DBName:     firstEnv([]string{"POSTGRES_DB", "DB_NAME"}, "alike"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		RedisAddr:     redisAddr(),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),

		JWTSecret:        getEnv("JWT_SECRET", defaultJWTSecret),
		JWTAccessExpire:  getEnvDuration("JWT_ACCESS_TTL", 120*time.Minute),
		JWTRefreshExpire: getEnvDuration("JWT_REFRESH_TTL", 168*time.Hour),

		MinIOEndpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinIOAccessKey: firstEnv([]string{"MINIO_ROOT_USER", "MINIO_ACCESS_KEY"}, "alike_minio_admin"),
		MinIOSecretKey: firstEnv([]string{"MINIO_ROOT_PASSWORD", "MINIO_SECRET_KEY"}, "minioadmin"),
		MinIOBucket:    getEnv("MINIO_BUCKET", "alike-uploads"),
		MinIOUseSSL:    getEnvBool("MINIO_USE_SSL", false),
		MinIOPublicURL: getEnv("MINIO_PUBLIC_URL", ""),

		UploadMaxImageBytes: int64(getEnvInt("UPLOAD_MAX_IMAGE_MB", 5)) << 20,
		UploadMaxDocBytes:   int64(getEnvInt("UPLOAD_MAX_DOC_MB", 10)) << 20,

		CORSAllowedOrigins: splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "")),
	}
}

// splitCSV 解析逗号分隔的配置值为去空白、去空项的字符串切片。
func splitCSV(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			out = append(out, v)
		}
	}
	return out
}

// IsProduction 报告是否运行于生产环境（APP_ENV=production）。
func (c *Config) IsProduction() bool {
	return strings.EqualFold(c.Env, "production")
}

// Validate 校验生产环境下的关键安全配置。
// 生产环境若仍使用默认 JWT 密钥则返回错误——否则攻击者可用公开于源码的密钥伪造任意用户令牌。
func (c *Config) Validate() error {
	if c.IsProduction() && c.JWTSecret == defaultJWTSecret {
		return errors.New("生产环境必须通过 JWT_SECRET 设置非默认密钥")
	}
	return nil
}

// redisAddr 优先用 REDIS_ADDR；否则由 REDIS_HOST + REDIS_PORT 组装。
func redisAddr() string {
	if v, ok := os.LookupEnv("REDIS_ADDR"); ok && v != "" {
		return v
	}
	host := getEnv("REDIS_HOST", "localhost")
	port := getEnv("REDIS_PORT", "6379")
	return host + ":" + port
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

// getEnvBool 解析布尔环境变量（true/1/yes/on 视为真），未设置或解析失败返回默认值。
func getEnvBool(key string, def bool) bool {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}

// firstEnv 依次尝试多个键，返回第一个非空值，全空则返回默认值。
func firstEnv(keys []string, def string) string {
	for _, k := range keys {
		if v, ok := os.LookupEnv(k); ok && v != "" {
			return v
		}
	}
	return def
}

// getEnvDuration 解析 Go duration 字符串（如 "15m"、"168h"），解析失败或未设置时返回默认值。
func getEnvDuration(key string, def time.Duration) time.Duration {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
