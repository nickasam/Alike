<script setup lang="ts">
/**
 * 登录页 —— Aurora 暗色玻璃拟态风格，居中卡片，不使用默认三列布局。
 * email + password，成功后跳转（支持 query.redirect 回跳，默认首页）。
 */
import { ApiError } from '~/composables/useApi'

definePageMeta({ layout: false, middleware: 'auth' })
useHead({ title: 'Alike · 登录' })

const auth = useAuth()
const route = useRoute()

const form = reactive({
  email: '',
  password: '',
})

const loading = ref(false)
const errorMsg = ref('')
const fieldErrors = reactive({
  email: '',
  password: '',
})

const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

function validate(): boolean {
  fieldErrors.email = ''
  fieldErrors.password = ''
  let ok = true

  if (!EMAIL_RE.test(form.email.trim())) {
    fieldErrors.email = '请输入有效的邮箱地址'
    ok = false
  }
  if (form.password.length < 6) {
    fieldErrors.password = '密码至少 6 位'
    ok = false
  }
  return ok
}

async function onSubmit() {
  errorMsg.value = ''
  if (!validate()) return

  loading.value = true
  try {
    await auth.login({ email: form.email.trim(), password: form.password })
    const redirect =
      typeof route.query.redirect === 'string' ? route.query.redirect : '/'
    await navigateTo(redirect)
  } catch (err) {
    errorMsg.value =
      err instanceof ApiError ? err.message : '登录失败，请稍后重试'
  } finally {
    loading.value = false
  }
}

const registerLink = computed(() => {
  const redirect =
    typeof route.query.redirect === 'string' ? route.query.redirect : undefined
  return redirect
    ? { path: '/register', query: { redirect } }
    : { path: '/register' }
})
</script>

<template>
  <div class="flex min-h-screen items-center justify-center px-4 py-10">
    <div class="glass-card animate-rise-in w-full max-w-[420px] p-8">
      <!-- 品牌头 -->
      <div class="mb-6 text-center">
        <h1 class="text-gradient text-2xl font-extrabold">Alike</h1>
        <p class="mt-1 text-sm text-dim">汇聚天下牛马，总有人懂你的辛苦</p>
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
            :aria-invalid="!!fieldErrors.email"
            class="rounded-md border border-border bg-surface-solid px-3 py-2 text-base text-text outline-none transition focus:border-ai-1"
          />
          <p v-if="fieldErrors.email" class="text-xs text-danger" role="alert">
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
            autocomplete="current-password"
            placeholder="至少 6 位"
            :aria-invalid="!!fieldErrors.password"
            class="rounded-md border border-border bg-surface-solid px-3 py-2 text-base text-text outline-none transition focus:border-ai-1"
          />
          <p v-if="fieldErrors.password" class="text-xs text-danger" role="alert">
            {{ fieldErrors.password }}
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
          {{ loading ? '处理中…' : '登录' }}
        </button>
      </form>

      <p class="mt-5 text-center text-sm text-mute">
        还没有账号？
        <NuxtLink :to="registerLink" class="text-ai-1 hover:underline">立即注册</NuxtLink>
      </p>
    </div>
  </div>
</template>
