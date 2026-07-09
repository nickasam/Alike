<script setup lang="ts">
/**
 * TopNav — 顶部导航（占位骨架）。
 * 玻璃拟态导航条 + Logo + 搜索 + 通知 + 主题切换 + 头像。
 * 移动端(<768)折叠为汉堡菜单入口（emit toggle-drawer）。
 */
defineEmits<{ (e: 'toggle-drawer'): void }>()

const { theme, toggle: toggleTheme, init: initTheme } = useTheme()
const auth = useAuth()
const router = useRouter()

/** 头像首字：用户昵称首字，回退为“牛”。 */
const avatarChar = computed(() => auth.user.value?.nickname?.[0] ?? '牛')
const menuOpen = ref(false)

async function onLogout() {
  menuOpen.value = false
  await auth.logout()
  await router.push('/')
}

onMounted(() => initTheme())
</script>

<template>
  <header
    class="sticky top-0 z-50 flex h-nav-h items-center gap-5 px-5 border-b border-border"
    style="
      background: var(--glass-bg);
      backdrop-filter: var(--glass-blur);
      -webkit-backdrop-filter: var(--glass-blur);
      box-shadow: var(--glass-shadow);
    "
  >
    <!-- 移动端汉堡菜单 -->
    <button
      class="grid h-10 w-10 place-items-center rounded-md border border-border text-dim md:hidden"
      aria-label="打开菜单"
      @click="$emit('toggle-drawer')"
    >
      <AppIcon name="menu" />
    </button>

    <!-- Logo -->
    <NuxtLink to="/" class="flex flex-shrink-0 items-center gap-2 font-bold">
      <span
        class="grid h-9 w-9 place-items-center rounded-md text-white bg-grad-ai shadow-glow-ai"
      >
        <AppIcon name="heart-handshake" />
      </span>
      <span class="text-gradient text-xl font-extrabold">Alike</span>
      <span class="hidden text-xs text-mute xl:inline">汇聚天下牛马</span>
    </NuxtLink>

    <!-- 搜索（桌面常驻） -->
    <div class="relative hidden max-w-[460px] flex-1 md:block">
      <span class="absolute left-3 top-1/2 -translate-y-1/2 text-mute">
        <AppIcon name="search" :size="18" />
      </span>
      <input
        type="search"
        aria-label="搜索"
        placeholder="搜索频道、日记、牛马..."
        class="h-10 w-full rounded-md border border-border bg-surface pl-10 pr-3 text-base text-text outline-none placeholder:text-mute focus:border-ai-1"
      />
    </div>

    <!-- 右侧操作 -->
    <div class="ml-auto flex items-center gap-3">
      <!-- 主题切换 -->
      <button
        class="grid h-10 w-10 place-items-center rounded-md border border-border bg-surface text-dim transition hover:text-ai-1"
        :aria-label="theme === 'dark' ? '切换到亮色' : '切换到暗色'"
        @click="toggleTheme"
      >
        <AppIcon :name="theme === 'dark' ? 'sun' : 'moon'" />
      </button>

      <!-- 通知 -->
      <button
        v-if="auth.isAuthenticated.value"
        class="relative hidden h-10 w-10 place-items-center rounded-md border border-border bg-surface text-dim transition hover:text-ai-1 sm:grid"
        aria-label="通知"
      >
        <AppIcon name="bell" />
        <span
          class="absolute -right-1 -top-1 grid h-[18px] min-w-[18px] place-items-center rounded-full bg-danger px-[5px] text-xs font-bold text-white"
          >3</span
        >
      </button>

      <!-- 已登录：头像 + 下拉菜单 -->
      <div v-if="auth.isAuthenticated.value" class="relative">
        <button
          class="flex items-center gap-2 rounded-md text-white"
          :aria-expanded="menuOpen"
          aria-haspopup="menu"
          aria-label="个人菜单"
          @click="menuOpen = !menuOpen"
        >
          <img
            v-if="auth.user.value?.avatar_url"
            :src="auth.user.value.avatar_url"
            :alt="auth.user.value?.nickname"
            class="h-10 w-10 rounded-md object-cover shadow-glow-ai"
          />
          <span
            v-else
            class="grid h-10 w-10 place-items-center rounded-md bg-grad-ai font-bold shadow-glow-ai"
          >
            {{ avatarChar }}
          </span>
          <span class="hidden text-sm font-medium text-text md:inline">
            {{ auth.user.value?.nickname }}
          </span>
        </button>

        <div
          v-if="menuOpen"
          role="menu"
          class="absolute right-0 top-12 z-50 w-44 overflow-hidden rounded-md border border-border bg-surface-solid shadow-lg"
        >
          <NuxtLink
            :to="`/profile/${auth.user.value?.id}`"
            role="menuitem"
            class="flex items-center gap-2 px-4 py-2.5 text-sm text-text transition hover:bg-surface"
            @click="menuOpen = false"
          >
            <AppIcon name="user" :size="18" />个人主页
          </NuxtLink>
          <button
            role="menuitem"
            class="flex w-full items-center gap-2 px-4 py-2.5 text-left text-sm text-danger transition hover:bg-surface"
            @click="onLogout"
          >
            <AppIcon name="log-out" :size="18" />退出登录
          </button>
        </div>
      </div>

      <!-- 未登录：登录 / 注册 -->
      <template v-else>
        <NuxtLink
          to="/login"
          class="rounded-md border border-border bg-surface px-3 py-2 text-sm font-medium text-dim transition hover:text-ai-1"
        >
          登录
        </NuxtLink>
        <NuxtLink
          to="/register"
          class="btn-primary rounded-md px-3 py-2 text-sm font-semibold"
        >
          注册
        </NuxtLink>
      </template>
    </div>
  </header>
</template>
