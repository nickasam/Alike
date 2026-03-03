import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/launcher',
    name: 'Launcher',
    component: () => import('@/views/Launcher.vue'),
    meta: { requiresAuth: false, title: 'Alike - 相似灵魂的相遇' }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: { requiresAuth: false, title: '登录 - Alike' }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/auth/Register.vue'),
    meta: { requiresAuth: false, title: '注册 - Alike' }
  },
  {
    path: '/',
    component: () => import('@/components/layout/MainLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('@/views/home/NearbyUsers.vue'),
        meta: { title: '附近用户 - Alike' }
      },
      {
        path: 'matches',
        name: 'Matches',
        component: () => import('@/views/match/Matches.vue'),
        meta: { title: '我的匹配 - Alike' }
      },
      {
        path: 'chat',
        name: 'ChatList',
        component: () => import('@/views/chat/ChatList.vue'),
        meta: { title: '聊天列表 - Alike' }
      },
      {
        path: 'chat/:userId',
        name: 'ChatRoom',
        component: () => import('@/views/chat/ChatRoom.vue'),
        meta: { title: '聊天 - Alike' }
      },
      {
        path: 'global',
        name: 'GlobalChat',
        component: () => import('@/views/global/GlobalChat.vue'),
        meta: { title: '全局聊天室 - Alike' }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/profile/Profile.vue'),
        meta: { title: '个人资料 - Alike' }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/auth/Login.vue'),
    meta: { title: '404 - 页面不存在' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  document.title = to.meta.title || 'Alike'
  
  // 检查是否需要认证
  const token = localStorage.getItem('alike_access_token')
  
  if (to.meta.requiresAuth && !token) {
    // 需要认证但没有 token，跳转到登录页
    next('/login')
  } else if ((to.path === '/login' || to.path === '/register') && token) {
    // 已登录用户访问登录页，跳转到首页
    next('/')
  } else {
    next()
  }
})

export default router