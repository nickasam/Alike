# ============================================================
# Alike 后端 Dockerfile — 多阶段构建
# 阶段一：golang 编译静态二进制
# 阶段二：alpine 精简运行镜像
# ============================================================

# ---------- build ----------
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 国内网络：使用七牛 GOPROXY 加速依赖下载
ENV GOPROXY=https://goproxy.cn,direct \
    GO111MODULE=on

# 先拷贝依赖清单，利用层缓存
COPY go.mod go.sum* ./
RUN go mod download

# 拷贝源码并编译（CGO 关闭，产出静态二进制）
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/server ./cmd/server

# ---------- runtime ----------
FROM alpine:3.20

# 时区 + CA 证书（HTTPS 调用 S3 等需要）
RUN apk add --no-cache ca-certificates tzdata wget && \
    adduser -D -u 10001 appuser

WORKDIR /app
COPY --from=builder /app/server /app/server
# 迁移文件随镜像一起，方便容器内执行迁移
COPY --from=builder /app/migrations /app/migrations

USER appuser

ENV BACKEND_PORT=8080
EXPOSE 8080

# 健康检查（后端需提供 /api/health 或 /healthz，任一即可）
HEALTHCHECK --interval=15s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- "http://127.0.0.1:${BACKEND_PORT}/api/health" || exit 1

ENTRYPOINT ["/app/server"]
