import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '~/stores/auth'

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

  it('setAuth 后写入 token 与 user，isAuthenticated 为 true', () => {
    const store = useAuthStore()
    store.setAuth('jwt-token', {
      id: 1,
      nickname: '牛马一号',
      level: 3,
      empathy_count: 42,
      anonymous: false,
    })
    expect(store.token).toBe('jwt-token')
    expect(store.user?.nickname).toBe('牛马一号')
    expect(store.isAuthenticated).toBe(true)
  })

  it('clear 后清空状态', () => {
    const store = useAuthStore()
    store.setAuth('jwt-token', {
      id: 1,
      nickname: '牛马一号',
      level: 3,
      empathy_count: 42,
      anonymous: false,
    })
    store.clear()
    expect(store.token).toBe('')
    expect(store.user).toBeNull()
    expect(store.isAuthenticated).toBe(false)
  })
})
