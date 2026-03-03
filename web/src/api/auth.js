import request from './request'

// 用户登录
export const login = (data) => {
  return request({
    url: '/auth/login',
    method: 'post',
    data
  })
}

// 用户注册
export const register = (data) => {
  return request({
    url: '/auth/register',
    method: 'post',
    data
  })
}

// 发送验证码
export const sendVerificationCode = (data) => {
  return request({
    url: '/auth/send-code',
    method: 'post',
    data
  })
}

// 刷新 token
export const refreshToken = (refreshToken) => {
  return request({
    url: '/auth/refresh',
    method: 'post',
    data: { refresh_token: refreshToken }
  })
}

// 登出
export const logout = () => {
  return request({
    url: '/auth/logout',
    method: 'post'
  })
}