<script setup lang="ts">
/**
 * 首页 — 热门频道 + 今日牛马榜（阶段一骨架）。
 * 体现 Aurora 暗色主题、渐变标题、玻璃卡片、响应式布局。
 */
useHead({ title: 'Alike · 汇聚天下牛马' })

const hotChannels = [
  { id: 1, name: '互联网大厂', desc: '996 的兄弟姐妹们，进来抱团', members: 1284 },
  { id: 3, name: '程序员', desc: 'bug 与 KPI 齐飞，头发共显示器一色', members: 986 },
  { id: 5, name: '摸鱼交流', desc: '带薪如厕经验分享大会', members: 2571 },
  { id: 6, name: '离职天堂', desc: '裸辞体验官招募中', members: 743 },
]
</script>

<template>
  <div class="flex flex-col gap-6">
    <!-- Hero -->
    <section class="glass-card animate-rise-in p-8">
      <h1 class="text-gradient text-2xl font-extrabold leading-tight md:text-3xl">
        汇聚天下牛马
      </h1>
      <p class="mt-2 text-md text-dim">
        总有人懂你的辛苦，你不是一个人在扛。
      </p>
      <div class="mt-5 flex flex-wrap gap-3">
        <button class="btn-primary px-5 py-2 text-base font-semibold">
          进入频道吐槽
        </button>
        <button
          class="rounded-md border border-border-strong px-5 py-2 text-base text-dim transition hover:text-text"
        >
          今日牛马榜
        </button>
      </div>
    </section>

    <!-- 热门频道 -->
    <section>
      <h2 class="mb-3 flex items-center gap-2 text-lg font-semibold">
        <AppIcon name="hash" :size="20" />
        热门频道
      </h2>
      <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
        <NuxtLink
          v-for="ch in hotChannels"
          :key="ch.id"
          :to="`/channel/${ch.id}`"
          class="glass-card block p-5"
        >
          <div class="flex items-center justify-between">
            <h3 class="text-md font-semibold">#{{ ch.name }}</h3>
            <span class="text-xs text-mute">{{ ch.members }} 位牛马</span>
          </div>
          <p class="mt-1 text-sm text-dim">{{ ch.desc }}</p>
        </NuxtLink>
      </div>
    </section>
  </div>
</template>
