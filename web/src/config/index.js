// 配置文件
export default {
  apiBaseUrl: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081/api/v1',
  appTitle: import.meta.env.VITE_APP_TITLE || 'Alike',
  appDescription: import.meta.env.VITE_APP_DESCRIPTION || '社交应用',
  
  // 本地存储键名
  storageKeys: {
    accessToken: 'alike_access_token',
    refreshToken: 'alike_refresh_token',
    userId: 'alike_user_id',
    userInfo: 'alike_user_info'
  },
  
  // 请求超时时间（毫秒）
  timeout: 10000
}