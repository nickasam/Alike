import request from './request'

// 获取全局聊天消息
export const getGlobalMessages = (params) => {
  return request({
    url: '/global/messages',
    method: 'get',
    params
  })
}

// 发送全局消息
export const sendGlobalMessage = (data) => {
  return request({
    url: '/global/messages',
    method: 'post',
    data
  })
}

// 获取在线用户列表
export const getOnlineUsers = () => {
  return request({
    url: '/global/online-users',
    method: 'get'
  })
}

// 获取在线用户数量
export const getOnlineCount = () => {
  return request({
    url: '/global/online-count',
    method: 'get'
  })
}
