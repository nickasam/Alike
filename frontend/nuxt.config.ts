// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-10',
  devtools: { enabled: true },

  modules: ['@nuxtjs/tailwindcss', '@pinia/nuxt'],

  // 组件按 PascalCase 文件名自动导入（不加目录前缀），
  // 使 TopNav / ChannelSidebar / EmotionBoard 等可直接使用。
  components: [{ path: '~/components', pathPrefix: false }],

  css: ['~/assets/css/main.css'],

  // 暗色 Aurora 极光风为默认主题，通过 <html class="dark"> 控制。
  app: {
    head: {
      htmlAttrs: { lang: 'zh-CN', class: 'dark' },
      title: 'Alike · 汇聚天下牛马',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        {
          name: 'description',
          content: 'Alike — 面向打工人的情感共鸣型聊天社区，总有人懂你的辛苦。',
        },
      ],
      link: [
        {
          rel: 'preconnect',
          href: 'https://fonts.googleapis.com',
        },
      ],
    },
  },

  runtimeConfig: {
    public: {
      // 后端 REST API 基址，可用环境变量 NUXT_PUBLIC_API_BASE 覆盖
      apiBase: process.env.NUXT_PUBLIC_API_BASE || '/api',
      // WebSocket 端点
      wsBase: process.env.NUXT_PUBLIC_WS_BASE || '/ws',
    },
  },

  typescript: {
    strict: true,
    typeCheck: false,
  },
})
