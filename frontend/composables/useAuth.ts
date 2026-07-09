/**
 * useAuth — 认证相关操作，桥接 auth store 与 useApi。
 * 阶段二：注册 / 登录 / 刷新 / 登出的完整流程。
 */
import { useAuthStore } from '~/stores/auth'
import type { AuthUser } from '~/stores/auth'

export interface LoginPayload {
  email: string
  password: string
}

export interface RegisterPayload {
  email: string
  password: string
  nickname: string
}

/** 后端 /auth/login 与 /auth/register 返回结构（对齐后端 AuthResponse） */
export interface AuthResult {
  user: AuthUser
  tokens: {
    access_token: string
    refresh_token: string
  }
}

export function useAuth() {
  const api = useApi()
  const store = useAuthStore()

  async function login(payload: LoginPayload) {
    const res = await api.post<AuthResult>('/auth/login', payload, {
      skipRefresh: true,
    })
    store.setAuth(res.tokens.access_token, res.tokens.refresh_token, res.user)
    return res
  }

  async function register(payload: RegisterPayload) {
    const res = await api.post<AuthResult>('/auth/register', payload, {
      skipRefresh: true,
    })
    store.setAuth(res.tokens.access_token, res.tokens.refresh_token, res.user)
    return res
  }

  /** 用 refresh token 换新 access token，写回 store。 */
  async function refreshToken() {
    const res = await api.post<{ tokens: { access_token: string } }>(
      '/auth/refresh',
      { refresh_token: store.refreshToken },
      { skipRefresh: true },
    )
    store.setToken(res.tokens.access_token)
    return res.tokens.access_token
  }

  async function fetchMe() {
    const user = await api.get<AuthUser>('/auth/me')
    store.setUser(user)
    return user
  }

  /** 登出：通知后端使 refresh token 失效，无论成功与否都清除本地态。 */
  async function logout() {
    try {
      await api.post('/auth/logout', { refresh_token: store.refreshToken }, {
        skipRefresh: true,
      })
    } catch {
      // 后端失败不阻塞本地登出
    } finally {
      store.clear()
    }
  }

  return {
    login,
    register,
    refreshToken,
    logout,
    fetchMe,
    isAuthenticated: computed(() => store.isAuthenticated),
    user: computed(() => store.user),
  }
}
