/**
 * useApi — 统一 REST API 封装。
 *
 * 后端响应统一格式：{ code, message, data }
 * 分页响应：{ code, message, data: { list, total, page, page_size } }
 *
 * 401（access token 过期）时自动用 refresh token 刷新并重试一次；
 * 刷新失败则清除登录态并跳转 /login。
 *
 * 用法：
 *   const api = useApi()
 *   const user = await api.get<User>('/users/1')
 *   const res = await api.post('/auth/login', { email, password })
 */
import { useAuthStore, TOKEN_KEY, REFRESH_KEY } from '~/stores/auth'

export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

export interface PageData<T = unknown> {
  list: T[]
  total: number
  page: number
  page_size: number
}

export class ApiError extends Error {
  code: number
  constructor(code: number, message: string) {
    super(message)
    this.name = 'ApiError'
    this.code = code
  }
}

type Query = Record<string, string | number | boolean | undefined>

interface RequestOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
  body?: unknown
  query?: Query
  /** 内部标记：已尝试过刷新，避免无限重试 */
  _retried?: boolean
  /** 内部标记：跳过 401 自动刷新（如刷新/登录本身） */
  _skipRefresh?: boolean
}

export function useApi() {
  const config = useRuntimeConfig()
  const baseURL = config.public.apiBase

  /** 读取当前 JWT（登录后写入）。SSR 阶段无 localStorage，返回空。 */
  function getToken(): string | null {
    if (import.meta.server) return null
    return localStorage.getItem(TOKEN_KEY)
  }

  function getRefreshToken(): string | null {
    if (import.meta.server) return null
    return localStorage.getItem(REFRESH_KEY)
  }

  /**
   * 用 refresh token 换取新的 access token。
   * 成功返回 true 并已写入 store/localStorage，失败返回 false。
   */
  async function tryRefresh(): Promise<boolean> {
    const refreshToken = getRefreshToken()
    if (!refreshToken) return false
    try {
      const res = await $fetch<ApiResponse<{ tokens: { access_token: string } }>>('/auth/refresh', {
        baseURL,
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: { refresh_token: refreshToken },
      })
      if (res && typeof res.code === 'number' && res.code !== 0) return false
      const token = res.data?.tokens?.access_token
      if (!token) return false
      useAuthStore().setToken(token)
      return true
    } catch {
      return false
    }
  }

  /** 刷新失败 / 无 refresh token：清除登录态并跳转登录页。 */
  function handleAuthFailure() {
    useAuthStore().clear()
    if (import.meta.client) {
      const route = useRoute()
      if (route.path !== '/login') {
        navigateTo({ path: '/login', query: { redirect: route.fullPath } })
      }
    }
  }

  function isUnauthorized(err: unknown): boolean {
    return (
      typeof err === 'object' &&
      err !== null &&
      (err as { statusCode?: number }).statusCode === 401
    )
  }

  async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
    const token = getToken()
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    }
    if (token) headers.Authorization = `Bearer ${token}`

    try {
      const res = await $fetch<ApiResponse<T>>(path, {
        baseURL,
        method: options.method ?? 'GET',
        headers,
        query: options.query,
        body: options.body as Record<string, unknown> | undefined,
      })

      // 统一 code 校验：code === 0 视为成功
      if (res && typeof res.code === 'number' && res.code !== 0) {
        throw new ApiError(res.code, res.message || '请求失败')
      }
      return res.data
    } catch (err) {
      // 401 且未跳过刷新、未重试过：尝试刷新后重试一次
      if (
        isUnauthorized(err) &&
        !options._skipRefresh &&
        !options._retried &&
        import.meta.client
      ) {
        const refreshed = await tryRefresh()
        if (refreshed) {
          return request<T>(path, { ...options, _retried: true })
        }
        handleAuthFailure()
      }
      throw err
    }
  }

  return {
    getToken,
    get: <T>(path: string, query?: Query) =>
      request<T>(path, { method: 'GET', query }),
    post: <T>(path: string, body?: unknown, opts?: { skipRefresh?: boolean }) =>
      request<T>(path, { method: 'POST', body, _skipRefresh: opts?.skipRefresh }),
    put: <T>(path: string, body?: unknown) =>
      request<T>(path, { method: 'PUT', body }),
    del: <T>(path: string, query?: Query) =>
      request<T>(path, { method: 'DELETE', query }),
  }
}
