<script setup lang="ts">
/**
 * 注册页 —— Aurora 暗色玻璃拟态风格，居中卡片。
 * email + password + nickname 必填，industry / job_title / work_years 选填。
 * 前端校验通过后调用 useAuth().register，成功自动登录并跳转（支持 redirect 回跳）。
 */
import { ApiError } from '~/composables/useApi'

definePageMeta({ layout: false, middleware: 'auth' })
useHead({ title: 'Alike · 注册' })

const auth = useAuth()
const route = useRoute()

const form = reactive({
  email: '',
  password: '',
  nickname: '',
  industry: '',
  job_title: '',
  work_years: '' as string,
})

const loading = ref(false)
const errorMsg = ref('')
const fieldErrors = reactive({
  email: '',
  password: '',
  nickname: '',
  work_years: '',
})

const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
void EMAIL_RE

function validate(): boolean {
  const errs = validateRegisterForm(form)
  fieldErrors.email = errs.email
  fieldErrors.password = errs.password
  fieldErrors.nickname = errs.nickname
  fieldErrors.work_years = errs.work_years
  return isFormValid(errs)
}

async function onSubmit() {
  errorMsg.value = ''
  if (!validate()) return

  loading.value = true
  try {
    await auth.register({
      email: form.email.trim(),
      password: form.password,
      nickname: form.nickname.trim(),
      industry: form.industry.trim() || undefined,
      job_title: form.job_title.trim() || undefined,
      work_years: form.work_years !== '' ? Number(form.work_years) : undefined,
    })
    const redirect =
      typeof route.query.redirect === 'string' ? route.query.redirect : '/'
    await navigateTo(redirect)
  } catch (err) {
    errorMsg.value =
      err instanceof ApiError ? err.message : '注册失败，请稍后重试'
  } finally {
    loading.value = false
  }
}

const loginLink = computed(() => {
  const redirect =
    typeof route.query.redirect === 'string' ? route.query.redirect : undefined
  return redirect ? { path: '/login', query: { redirect } } : { path: '/login' }
})
</script>

<template>
  <div class="flex min-h-screen items-center justify-center px-4 py-10">
    <div class="glass-card animate-rise-in w-full max-w-[420px] p-8">
      <!-- 品牌头 -->
      <div class="mb-6 text-center">
        <h1 class="text-gradient text-2xl font-extrabold">Alike</h1>
        <p class="mt-1 text-sm text-dim">加入我们，总有人懂你的辛苦</p>
      </div>

      <form class="flex flex-col gap-4" novalidate @submit.prevent="onSubmit">
        <!-- 邮箱 -->
        <div class="flex flex-col gap-1">
          <label for="email" class="text-sm text-dim">邮箱 <span class="text-danger">*</span></label>
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
          <label for="password" class="text-sm text-dim">密码 <span class="text-danger">*</span></label>
          <input
            id="password"
            v-model="form.password"
            type="password"
            autocomplete="new-password"
            placeholder="至少 6 位"
            :aria-invalid="!!fieldErrors.password"
            class="rounded-md border border-border bg-surface-solid px-3 py-2 text-base text-text outline-none transition focus:border-ai-1"
          />
          <p v-if="fieldErrors.password" class="text-xs text-danger" role="alert">
            {{ fieldErrors.password }}
          </p>
        </div>

        <!-- 昵称 -->
        <div class="flex flex-col gap-1">
          <label for="nickname" class="text-sm text-dim">昵称 <span class="text-danger">*</span></label>
          <input
            id="nickname"
            v-model="form.nickname"
            type="text"
            autocomplete="nickname"
            placeholder="给自己起个牛马名"
            :aria-invalid="!!fieldErrors.nickname"
            class="rounded-md border border-border bg-surface-solid px-3 py-2 text-base text-text outline-none transition focus:border-ai-1"
          />
          <p v-if="fieldErrors.nickname" class="text-xs text-danger" role="alert">
            {{ fieldErrors.nickname }}
          </p>
        </div>

        <!-- 选填字段 -->
        <div class="grid grid-cols-2 gap-3">
          <div class="flex flex-col gap-1">
            <label for="industry" class="text-sm text-dim">行业</label>
            <input
              id="industry"
              v-model="form.industry"
              type="text"
              placeholder="选填"
              class="rounded-md border border-border bg-surface-solid px-3 py-2 text-base text-text outline-none transition focus:border-ai-1"
            />
          </div>
          <div class="flex flex-col gap-1">
            <label for="job_title" class="text-sm text-dim">岗位</label>
            <input
              id="job_title"
              v-model="form.job_title"
              type="text"
              placeholder="选填"
              class="rounded-md border border-border bg-surface-solid px-3 py-2 text-base text-text outline-none transition focus:border-ai-1"
            />
          </div>
        </div>

        <div class="flex flex-col gap-1">
          <label for="work_years" class="text-sm text-dim">工龄（年）</label>
          <input
            id="work_years"
            v-model="form.work_years"
            type="number"
            min="0"
            max="60"
            placeholder="选填"
            :aria-invalid="!!fieldErrors.work_years"
            class="rounded-md border border-border bg-surface-solid px-3 py-2 text-base text-text outline-none transition focus:border-ai-1"
          />
          <p v-if="fieldErrors.work_years" class="text-xs text-danger" role="alert">
            {{ fieldErrors.work_years }}
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
          {{ loading ? '处理中…' : '注册' }}
        </button>
      </form>

      <p class="mt-5 text-center text-sm text-mute">
        已有账号？
        <NuxtLink :to="loginLink" class="text-ai-1 hover:underline">去登录</NuxtLink>
      </p>
    </div>
  </div>
</template>
