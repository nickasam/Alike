import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mockNuxtImport } from '@nuxt/test-utils/runtime'
import authMiddleware from '~/middleware/auth'

// 回归测试：/register 曾因未被视为公开认证页而被 auth 中间件
// 反复跳回 /login?redirect=/register，导致注册页无法访问（无限回跳）。

let authed = false
const { navigateSpy } = vi.hoisted(() => ({ navigateSpy: vi.fn((arg: unknown) => arg) }))

mockNuxtImport('navigateTo', () => navigateSpy)
mockNuxtImport('useAuthStore', () => () => ({
  get isAuthenticated() {
    return authed
  },
}))

function makeTo(path: string) {
  return { path, fullPath: path } as any
}

describe('auth route middleware', () => {
  beforeEach(() => {
    authed = false
    navigateSpy.mockClear()
  })

  it('未登录访问 /register 不应被拦截回跳（回归）', () => {
    const result = authMiddleware(makeTo('/register'), makeTo('/'))
    expect(navigateSpy).not.toHaveBeenCalled()
    expect(result).toBeUndefined()
  })

  it('未登录访问 /login 不应被拦截', () => {
    const result = authMiddleware(makeTo('/login'), makeTo('/'))
    expect(navigateSpy).not.toHaveBeenCalled()
    expect(result).toBeUndefined()
  })

  it('未登录访问受保护页跳转 /login 并带 redirect', () => {
    authMiddleware(makeTo('/diary'), makeTo('/'))
    expect(navigateSpy).toHaveBeenCalledWith({
      path: '/login',
      query: { redirect: '/diary' },
    })
  })

  it('已登录访问 /register 跳回首页', () => {
    authed = true
    authMiddleware(makeTo('/register'), makeTo('/'))
    expect(navigateSpy).toHaveBeenCalledWith('/')
  })

  it('已登录访问 /login 跳回首页', () => {
    authed = true
    authMiddleware(makeTo('/login'), makeTo('/'))
    expect(navigateSpy).toHaveBeenCalledWith('/')
  })
})
