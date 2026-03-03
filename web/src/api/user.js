import request from './request'

// 获取附近用户
export const getNearbyUsers = (params) => {
  return request({
    url: '/users/nearby',
    method: 'get',
    params
  })
}

// 获取用户信息
export const getUserInfo = (userId) => {
  return request({
    url: `/users/${userId}`,
    method: 'get'
  })
}

// 更新用户资料
export const updateProfile = (data) => {
  return request({
    url: '/users/profile',
    method: 'put',
    data
  })
}

// 上传头像
export const uploadAvatar = (formData) => {
  return request({
    url: '/users/avatar',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 喜欢用户
export const likeUser = (userId) => {
  return request({
    url: `/users/${userId}/like`,
    method: 'post'
  })
}

// 获取匹配列表
export const getMatches = () => {
  return request({
    url: '/matches',
    method: 'get'
  })
}