package config

import (
	"testing"
	"time"
)

// TestLoadComposeContract 验证配置读取与 docker-compose/.env 契约的变量名一致。
func TestLoadComposeContract(t *testing.T) {
	envs := map[string]string{
		"APP_ENV":         "production",
		"BACKEND_PORT":    "8080",
		"POSTGRES_HOST":   "postgres",
		"POSTGRES_PORT":   "5432",
		"POSTGRES_USER":   "alike",
		"POSTGRES_DB":     "alike",
		"REDIS_HOST":      "redis",
		"REDIS_PORT":      "6379",
		"JWT_ACCESS_TTL":  "15m",
		"JWT_REFRESH_TTL": "168h",
	}
	for k, v := range envs {
		t.Setenv(k, v)
	}

	c := Load()

	if c.DBHost != "postgres" {
		t.Errorf("DBHost = %q, want postgres", c.DBHost)
	}
	if c.ServerPort != "8080" {
		t.Errorf("ServerPort = %q, want 8080", c.ServerPort)
	}
	if c.RedisAddr != "redis:6379" {
		t.Errorf("RedisAddr = %q, want redis:6379", c.RedisAddr)
	}
	if c.JWTAccessExpire != 15*time.Minute {
		t.Errorf("JWTAccessExpire = %v, want 15m", c.JWTAccessExpire)
	}
	if c.JWTRefreshExpire != 168*time.Hour {
		t.Errorf("JWTRefreshExpire = %v, want 168h", c.JWTRefreshExpire)
	}
}

// TestLoadLegacyFallback 验证旧命名（DB_*/REDIS_ADDR/SERVER_PORT）仍可回退生效。
func TestLoadLegacyFallback(t *testing.T) {
	t.Setenv("DB_HOST", "legacy-db")
	t.Setenv("SERVER_PORT", "9090")
	t.Setenv("REDIS_ADDR", "legacy-redis:6380")

	c := Load()

	if c.DBHost != "legacy-db" {
		t.Errorf("DBHost = %q, want legacy-db", c.DBHost)
	}
	if c.ServerPort != "9090" {
		t.Errorf("ServerPort = %q, want 9090", c.ServerPort)
	}
	if c.RedisAddr != "legacy-redis:6380" {
		t.Errorf("RedisAddr = %q, want legacy-redis:6380", c.RedisAddr)
	}
}

// TestLoadDefaults 验证无环境变量时的默认值。
func TestLoadDefaults(t *testing.T) {
	c := Load()

	if c.DBHost != "localhost" {
		t.Errorf("DBHost = %q, want localhost", c.DBHost)
	}
	if c.RedisAddr != "localhost:6379" {
		t.Errorf("RedisAddr = %q, want localhost:6379", c.RedisAddr)
	}
	if c.JWTAccessExpire != 120*time.Minute {
		t.Errorf("JWTAccessExpire = %v, want 120m", c.JWTAccessExpire)
	}
}
