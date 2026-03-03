import request from './request'

// 获取聊天列表
export const getChatList = () => {
  return request({
    url: '/chats',
    method: 'get'
  })
}

// 获取与特定用户的聊天记录
export const getChatMessages = (userId, params) => {
  return request({
    url: `/chats/${userId}/messages`,
    method: 'get',
    params
  })
}

// 发送消息
export const sendMessage = (userId, data) => {
  return request({
    url: `/chats/${userId}/messages`,
    method: 'post',
    data
  })
}

// 获取未读消息数量
export const getUnreadCount = () => {
  return request({
    url: '/chats/unread-count',
    method: 'get'
  })
}

// 标记消息为已读
export const markAsRead = (userId) => {
  return request({
    url: `/chats/${userId}/read`,
    method: 'post'
  })
}