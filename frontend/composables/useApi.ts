/**
 * useApi — 统一 REST API 封装。
 *
 * 后端响应统一格式：{ code, message, data }
 * 分页响应：{ code, message, data: { list, total, page, page_size } }
 *
 * 用法：
 *   const api = useApi()
 *   const user = await api.get<User>('/users/1')
 *   const res = await api.post('/auth/login', { email, password })
 */

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

export function useApi() {
  const config = useRuntimeConfig()
  const baseURL = config.public.apiBase

  /** 读取当前 JWT（登录后写入）。SSR 阶段无 localStorage，返回空。 */
  function getToken(): string | null {
    if (import.meta.server) return null
    return localStorage.getItem('alike_token')
  }

  async function request<T>(
    path: string,
    options: {
      method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
      body?: unknown
      query?: Query
    } = {},
  ): Promise<T> {
    const token = getToken()
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    }
    if (token) headers.Authorization = `Bearer ${token}`

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
  }

  return {
    getToken,
    get: <T>(path: string, query?: Query) =>
      request<T>(path, { method: 'GET', query }),
    post: <T>(path: string, body?: unknown) =>
      request<T>(path, { method: 'POST', body }),
    put: <T>(path: string, body?: unknown) =>
      request<T>(path, { method: 'PUT', body }),
    del: <T>(path: string, query?: Query) =>
      request<T>(path, { method: 'DELETE', query }),
  }
}
