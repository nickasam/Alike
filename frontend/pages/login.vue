<script setup lang="ts">
/**
 * 登录 / 注册页 —— 同一页面 Tab 切换。
 * Aurora 暗色玻璃拟态风格，居中卡片，不使用默认三列布局。
 */
import { ApiError } from '~/composables/useApi'

definePageMeta({ layout: false, middleware: 'auth' })
useHead({ title: 'Alike · 登录 / 注册' })

const auth = useAuth()
const route = useRoute()

type Tab = 'login' | 'register'
const tab = ref<Tab>('login')

const form = reactive({
  email: '',
  password: '',
  nickname: '',
})

const loading = ref(false)
const errorMsg = ref('')
const fieldErrors = reactive({
  email: '',
  password: '',
  nickname: '',
})

const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

function switchTab(next: Tab) {
  if (tab.value === next) return
  tab.value = next
  errorMsg.value = ''
  fieldErrors.email = ''
  fieldErrors.password = ''
  fieldErrors.nickname = ''
}

function validate(): boolean {
  fieldErrors.email = ''
  fieldErrors.password = ''
  fieldErrors.nickname = ''
  let ok = true

  if (!EMAIL_RE.test(form.email.trim())) {
    fieldErrors.email = '请输入有效的邮箱地址'
    ok = false
  }
  if (form.password.length < 6) {
    fieldErrors.password = '密码至少 6 位'
    ok = false
  }
  if (tab.value === 'register' && !form.nickname.trim()) {
    fieldErrors.nickname = '昵称不能为空'
    ok = false
  }
  return ok
}

async function onSubmit() {
  errorMsg.value = ''
  if (!validate()) return

  loading.value = true
  try {
    if (tab.value === 'login') {
      await auth.login({ email: form.email.trim(), password: form.password })
    } else {
      await auth.register({
        email: form.email.trim(),
        password: form.password,
        nickname: form.nickname.trim(),
      })
    }
    const redirect =
      typeof route.query.redirect === 'string' ? route.query.redirect : '/'
    await navigateTo(redirect)
  } catch (err) {
    errorMsg.value =
      err instanceof ApiError
        ? err.message
        : tab.value === 'login'
          ? '登录失败，请稍后重试'
          : '注册失败，请稍后重试'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="flex min-h-screen items-center justify-center px-4 py-10">
    <div class="glass-card animate-rise-in w-full max-w-[420px] p-8">
      <!-- 品牌头 -->
      <div class="mb-6 text-center">
        <h1 class="text-gradient text-2xl font-extrabold">Alike</h1>
        <p class="mt-1 text-sm text-dim">汇聚天下牛马，总有人懂你的辛苦</p>
      </div>

      <!-- Tab 切换 -->
      <div
        class="mb-6 flex rounded-md border border-border p-1"
        role="tablist"
      >
        <button
          type="button"
          role="tab"
          :aria-selected="tab === 'login'"
          class="flex-1 rounded-sm py-2 text-base font-medium transition"
          :class="
            tab === 'login'
              ? 'bg-grad-ai text-white shadow-glow-ai'
              : 'text-dim hover:text-text'
          "
          @click="switchTab('login')"
        >
          登录
        </button>
        <button
          type="button"
          role="tab"
          :aria-selected="tab === 'register'"
          class="flex-1 rounded-sm py-2 text-base font-medium transition"
          :class="
            tab === 'register'
              ? 'bg-grad-ai text-white shadow-glow-ai'
              : 'text-dim hover:text-text'
          "
          @click="switchTab('register')"
        >
          注册
        </button>
      </div>

      <form class="flex flex-col gap-4" novalidate @submit.prevent="onSubmit">
        <!-- 邮箱 -->
        <div class="flex flex-col gap-1">
          <label for="email" class="text-sm text-dim">邮箱</label>
          <input
            id="email"
            v-model="form.email"
            type="email"
            autocomplete="email"
            placeholder="you@example.com"
            class="rounded-md border border-border bg-surface-solid px-3 py-2 text-base text-text outline-none transition focus:border-ai-1"
          />
          <p v-if="fieldErrors.email" class="text-xs text-danger">
            {{ fieldErrors.email }}
          </p>
        </div>

        <!-- 密码 -->
        <div class="flex flex-col gap-1">
          <label for="password" class="text-sm text-dim">密码</label>
          <input
            id="password"
            v-model="form.password"
            type="password"
            :autocomplete="tab === 'login' ? 'current-password' : 'new-password'"
            placeholder="至少 6 位"
            class="rounded-md border border-border bg-surface-solid px-3 py-2 text-base text-text outline-none transition focus:border-ai-1"
          />
          <p v-if="fieldErrors.password" class="text-xs text-danger">
            {{ fieldErrors.password }}
          </p>
        </div>

        <!-- 昵称（仅注册） -->
        <div v-if="tab === 'register'" class="flex flex-col gap-1">
          <label for="nickname" class="text-sm text-dim">昵称</label>
          <input
            id="nickname"
            v-model="form.nickname"
            type="text"
            autocomplete="nickname"
            placeholder="给自己起个牛马名"
            class="rounded-md border border-border bg-surface-solid px-3 py-2 text-base text-text outline-none transition focus:border-ai-1"
          />
          <p v-if="fieldErrors.nickname" class="text-xs text-danger">
            {{ fieldErrors.nickname }}
          </p>
        </div>

        <!-- 全局错误 -->
        <p
          v-if="errorMsg"
          class="rounded-md border border-danger/40 bg-danger/10 px-3 py-2 text-sm text-danger"
          role="alert"
        >
          {{ errorMsg }}
        </p>

        <button
          type="submit"
          :disabled="loading"
          class="btn-primary mt-2 py-2.5 text-base font-semibold disabled:cursor-not-allowed disabled:opacity-60"
        >
          {{ loading ? '处理中…' : tab === 'login' ? '登录' : '注册' }}
        </button>
      </form>

      <p class="mt-5 text-center text-sm text-mute">
        {{ tab === 'login' ? '还没有账号？' : '已有账号？' }}
        <button
          type="button"
          class="text-ai-1 hover:underline"
          @click="switchTab(tab === 'login' ? 'register' : 'login')"
        >
          {{ tab === 'login' ? '立即注册' : '去登录' }}
        </button>
      </p>
    </div>
  </div>
</template>
