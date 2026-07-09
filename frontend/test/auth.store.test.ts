import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '~/stores/auth'

const sampleUser = {
  id: 1,
  nickname: '牛马一号',
  level: 3,
  empathy_count: 42,
  anonymous: false,
}

describe('auth store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    // 稳定的内存版 localStorage，避免测试环境实现差异
    const mem: Record<string, string> = {}
    vi.stubGlobal('localStorage', {
      getItem: (k: string) => mem[k] ?? null,
      setItem: (k: string, v: string) => {
        mem[k] = v
      },
      removeItem: (k: string) => {
        delete mem[k]
      },
    })
  })

  it('未登录时 isAuthenticated 为 false', () => {
    const store = useAuthStore()
    expect(store.isAuthenticated).toBe(false)
  })

  it('setAuth 后写入 token/refreshToken/user，isAuthenticated 为 true', () => {
    const store = useAuthStore()
    store.setAuth('jwt-token', 'refresh-token', sampleUser)
    expect(store.token).toBe('jwt-token')
    expect(store.refreshToken).toBe('refresh-token')
    expect(store.user?.nickname).toBe('牛马一号')
    expect(store.isAuthenticated).toBe(true)
  })

  it('setToken 仅更新 access token', () => {
    const store = useAuthStore()
    store.setAuth('jwt-token', 'refresh-token', sampleUser)
    store.setToken('new-jwt')
    expect(store.token).toBe('new-jwt')
    expect(store.refreshToken).toBe('refresh-token')
  })

  it('init 从 localStorage 恢复登录态', () => {
    localStorage.setItem('alike_token', 'saved-token')
    localStorage.setItem('alike_refresh_token', 'saved-refresh')
    localStorage.setItem('alike_user', JSON.stringify(sampleUser))
    const store = useAuthStore()
    store.init()
    expect(store.token).toBe('saved-token')
    expect(store.refreshToken).toBe('saved-refresh')
    expect(store.user?.nickname).toBe('牛马一号')
    expect(store.isAuthenticated).toBe(true)
  })

  it('clear 后清空状态与 localStorage', () => {
    const store = useAuthStore()
    store.setAuth('jwt-token', 'refresh-token', sampleUser)
    store.clear()
    expect(store.token).toBe('')
    expect(store.refreshToken).toBe('')
    expect(store.user).toBeNull()
    expect(store.isAuthenticated).toBe(false)
    expect(localStorage.getItem('alike_token')).toBeNull()
    expect(localStorage.getItem('alike_refresh_token')).toBeNull()
  })
})

