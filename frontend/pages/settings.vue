<script setup lang="ts">
/**
 * 个人资料编辑页 — 修改昵称/简介/行业/岗位/工龄/匿名开关 + 头像上传。
 * 保存经 PUT /users/:id；头像先 POST /api/upload 拿 URL 再回填 avatar_url。
 * 成功后同步 authStore，顶栏头像/昵称即时更新。
 */
import { ApiError } from '~/composables/useApi'
import { useAuthStore, type AuthUser } from '~/stores/auth'

definePageMeta({ middleware: 'auth' })
useHead({ title: '编辑资料 · Alike' })

const api = useApi()
const authStore = useAuthStore()
const config = useRuntimeConfig()

// 用当前登录用户预填表单。
const form = reactive({
  nickname: '',
  bio: '',
  industry: '',
  job_title: '',
  work_years: 0,
  is_anonymous: false,
  avatar_url: '',
})

onMounted(() => {
  const u = authStore.user
  if (!u) return
  form.nickname = u.nickname ?? ''
  form.bio = u.bio ?? ''
  form.industry = u.industry ?? ''
  form.job_title = u.job_title ?? ''
  form.work_years = u.work_years ?? 0
  form.is_anonymous = u.is_anonymous ?? false
  form.avatar_url = u.avatar_url ?? ''
})

const loading = ref(false)
const errorMsg = ref('')
const okMsg = ref('')
const fieldErrors = reactive<Record<string, string>>({})

// —— 头像上传 —— //
const uploading = ref(false)
const avatarError = ref('')

async function onPickAvatar(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  avatarError.value = ''
  if (!file.type.startsWith('image/')) {
    avatarError.value = '请选择图片文件'
    return
  }
  if (file.size > 5 * 1024 * 1024) {
    avatarError.value = '图片不能超过 5MB'
    return
  }
  uploading.value = true
  try {
    const fd = new FormData()
    fd.append('file', file)
    // useApi 只支持 JSON，这里直接用 $fetch 传 multipart（不手动设 Content-Type）。
    const token = api.getToken()
    const res = await $fetch<{ code: number; message: string; data: { url: string } }>(
      `${config.public.apiBase}/upload`,
      {
        method: 'POST',
        body: fd,
        headers: token ? { Authorization: `Bearer ${token}` } : {},
      },
    )
    if (res.code !== 0) {
      avatarError.value = res.message || '上传失败'
      return
    }
    form.avatar_url = res.data.url
  } catch (err: any) {
    avatarError.value = err?.data?.message || '上传失败，请稍后重试'
  } finally {
    uploading.value = false
    input.value = '' // 允许重复选择同一文件
  }
}

function validate(): boolean {
  Object.keys(fieldErrors).forEach((k) => delete fieldErrors[k])
  if (!form.nickname.trim()) fieldErrors.nickname = '昵称不能为空'
  else if (form.nickname.length > 100) fieldErrors.nickname = '昵称最多 100 字'
  if (form.bio.length > 200) fieldErrors.bio = '简介最多 200 字'
  if (form.work_years < 0) fieldErrors.work_years = '工龄不能为负'
  return Object.keys(fieldErrors).length === 0
}

async function onSubmit() {
  errorMsg.value = ''
  okMsg.value = ''
  if (!validate()) return
  const uid = authStore.user?.id
  if (!uid) {
    errorMsg.value = '未登录'
    return
  }
  loading.value = true
  try {
    const updated = await api.put<AuthUser>(`/users/${uid}`, {
      nickname: form.nickname.trim(),
      bio: form.bio,
      industry: form.industry,
      job_title: form.job_title,
      work_years: form.work_years,
      is_anonymous: form.is_anonymous,
      avatar_url: form.avatar_url,
    })
    authStore.setUser(updated)
    okMsg.value = '资料已保存'
  } catch (err) {
    errorMsg.value = err instanceof ApiError ? err.message : '保存失败，请稍后重试'
  } finally {
    loading.value = false
  }
}

const avatarChar = computed(() => form.nickname?.[0] ?? '牛')
</script>

<template>
  <div class="mx-auto w-full max-w-xl">
    <section class="glass-card animate-rise-in p-8">
      <h1 class="mb-6 flex items-center gap-2 text-xl font-semibold">
        <AppIcon name="user" :size="22" />
        编辑资料
      </h1>

      <!-- 头像 -->
      <div class="mb-6 flex items-center gap-4">
        <img
          v-if="form.avatar_url"
          :src="form.avatar_url"
          alt="头像"
          class="h-16 w-16 rounded-lg object-cover"
        />
        <span
          v-else
          class="grid h-16 w-16 place-items-center rounded-lg bg-grad-ai text-xl font-semibold text-white"
        >
          {{ avatarChar }}
        </span>
        <div class="flex flex-col gap-1">
          <label
            class="inline-flex w-fit cursor-pointer items-center gap-2 rounded-md border border-border-strong px-3 py-1.5 text-sm text-dim transition hover:text-text"
          >
            <AppIcon name="plus" :size="16" />
            {{ uploading ? '上传中…' : '更换头像' }}
            <input
              type="file"
              accept="image/*"
              class="hidden"
              :disabled="uploading"
              @change="onPickAvatar"
            />
          </label>
          <p v-if="avatarError" class="text-xs text-danger" role="alert">{{ avatarError }}</p>
          <p v-else class="text-xs text-mute">JPG/PNG/GIF/WebP，≤5MB</p>
        </div>
      </div>

      <form class="flex flex-col gap-4" @submit.prevent="onSubmit">
        <div class="flex flex-col gap-1">
          <label for="nickname" class="text-sm text-dim">昵称</label>
          <input
            id="nickname"
            v-model="form.nickname"
            type="text"
            maxlength="100"
            class="rounded-md border border-border bg-surface-solid px-3 py-2 text-base text-text outline-none transition focus:border-ai-1"
          />
          <p v-if="fieldErrors.nickname" class="text-xs text-danger" role="alert">
            {{ fieldErrors.nickname }}
          </p>
        </div>

        <div class="flex flex-col gap-1">
          <label for="bio" class="text-sm text-dim">简介</label>
          <textarea
            id="bio"
            v-model="form.bio"
            rows="3"
            maxlength="200"
            placeholder="说点什么介绍自己…"
            class="resize-none rounded-md border border-border bg-surface-solid px-3 py-2 text-base text-text outline-none transition focus:border-ai-1"
          />
          <p class="text-right text-xs text-mute">{{ form.bio.length }}/200</p>
          <p v-if="fieldErrors.bio" class="text-xs text-danger" role="alert">
            {{ fieldErrors.bio }}
          </p>
        </div>

        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div class="flex flex-col gap-1">
            <label for="industry" class="text-sm text-dim">行业</label>
            <input
              id="industry"
              v-model="form.industry"
              type="text"
              maxlength="100"
              placeholder="如：互联网"
              class="rounded-md border border-border bg-surface-solid px-3 py-2 text-base text-text outline-none transition focus:border-ai-1"
            />
          </div>
          <div class="flex flex-col gap-1">
            <label for="job_title" class="text-sm text-dim">岗位</label>
            <input
              id="job_title"
              v-model="form.job_title"
              type="text"
              maxlength="100"
              placeholder="如：程序员"
              class="rounded-md border border-border bg-surface-solid px-3 py-2 text-base text-text outline-none transition focus:border-ai-1"
            />
          </div>
        </div>

        <div class="flex flex-col gap-1">
          <label for="work_years" class="text-sm text-dim">工龄（年）</label>
          <input
            id="work_years"
            v-model.number="form.work_years"
            type="number"
            min="0"
            class="w-32 rounded-md border border-border bg-surface-solid px-3 py-2 text-base text-text outline-none transition focus:border-ai-1"
          />
          <p v-if="fieldErrors.work_years" class="text-xs text-danger" role="alert">
            {{ fieldErrors.work_years }}
          </p>
        </div>

        <label class="flex items-center gap-2 text-sm text-dim">
          <input v-model="form.is_anonymous" type="checkbox" class="h-4 w-4 accent-ai-1" />
          默认匿名发言
        </label>

        <p
          v-if="errorMsg"
          class="rounded-md border border-danger/40 bg-danger/10 px-3 py-2 text-sm text-danger"
          role="alert"
        >
          {{ errorMsg }}
        </p>
        <p
          v-if="okMsg"
          class="rounded-md border border-empathy/40 bg-empathy/10 px-3 py-2 text-sm text-empathy"
          role="status"
        >
          {{ okMsg }}
        </p>

        <button
          type="submit"
          :disabled="loading"
          class="btn-primary mt-2 py-2.5 text-base font-semibold disabled:cursor-not-allowed disabled:opacity-60"
        >
          {{ loading ? '保存中…' : '保存' }}
        </button>
      </form>
    </section>
  </div>
</template>
