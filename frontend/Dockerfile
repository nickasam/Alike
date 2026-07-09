# ============================================================
# Alike 前端 Dockerfile — 多阶段构建
# 阶段一：node 安装依赖 + nuxt build
# 阶段二：node 精简镜像，运行 .output/server (Nitro)
# ============================================================

# ---------- build ----------
FROM node:20-alpine AS builder

WORKDIR /app

# 先拷贝依赖清单，利用层缓存（兼容 npm / pnpm / yarn 锁文件）
COPY package.json package-lock.json* pnpm-lock.yaml* yarn.lock* ./
RUN if [ -f pnpm-lock.yaml ]; then \
        corepack enable && pnpm install --frozen-lockfile; \
    elif [ -f yarn.lock ]; then \
        yarn install --frozen-lockfile; \
    else \
        npm install; \
    fi

# 拷贝源码并构建
COPY . .
RUN npm run build

# ---------- runtime ----------
FROM node:20-alpine AS runner

WORKDIR /app

RUN addgroup -g 10001 nodegrp && \
    adduser -D -u 10001 -G nodegrp nuxtuser

ENV NODE_ENV=production \
    NITRO_HOST=0.0.0.0 \
    NITRO_PORT=3000

# 仅拷贝 Nuxt 产物，运行期无需 node_modules（Nitro 已打包）
COPY --from=builder /app/.output /app/.output

USER nuxtuser

EXPOSE 3000

HEALTHCHECK --interval=15s --timeout=5s --start-period=15s --retries=3 \
    CMD wget -qO- "http://127.0.0.1:${NITRO_PORT}/" || exit 1

CMD ["node", ".output/server/index.mjs"]
