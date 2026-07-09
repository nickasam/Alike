/**
 * useAuth — 认证相关操作骨架，桥接 auth store 与 useApi。
 * 阶段二完善注册/登录/刷新逻辑，这里只提供最小可用形状。
 */
import { useAuthStore } from '~/stores/auth'

export interface LoginPayload {
  email: string
  password: string
}

export interface AuthResult {
  token: string
  user: import('~/stores/auth').AuthUser
}

export function useAuth() {
  const api = useApi()
  const store = useAuthStore()

  async function login(payload: LoginPayload) {
    const res = await api.post<AuthResult>('/auth/login', payload)
    store.setAuth(res.token, res.user)
    return res
  }

  async function fetchMe() {
    const user = await api.get<import('~/stores/auth').AuthUser>('/auth/me')
    store.setUser(user)
    return user
  }

  function logout() {
    store.clear()
  }

  return {
    login,
    logout,
    fetchMe,
    isAuthenticated: computed(() => store.isAuthenticated),
    user: computed(() => store.user),
  }
}
