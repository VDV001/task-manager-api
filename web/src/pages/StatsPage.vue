<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useTaskStore } from '@/stores/tasks'
import { ListTodo, Clock, Loader, CheckCircle2, AlertTriangle } from 'lucide-vue-next'

const { t } = useI18n()
const taskStore = useTaskStore()

onMounted(() => {
  taskStore.fetchStats()
})

const cards = computed(() => {
  const s = taskStore.stats
  if (!s) return []
  return [
    {
      labelKey: 'stats.total',
      value: s.total,
      icon: ListTodo,
      color: 'text-accent',
      bg: 'bg-accent/10',
      borderGlow: 'hover:border-accent/30',
      shadowGlow: 'hover:shadow-[0_0_30px_rgba(99,102,241,0.15)]',
    },
    {
      labelKey: 'stats.new',
      value: s.by_status.new ?? 0,
      icon: Clock,
      color: 'text-info',
      bg: 'bg-info/10',
      borderGlow: 'hover:border-info/30',
      shadowGlow: 'hover:shadow-[0_0_30px_rgba(56,189,248,0.15)]',
    },
    {
      labelKey: 'stats.inProgress',
      value: s.by_status.in_progress ?? 0,
      icon: Loader,
      color: 'text-warning',
      bg: 'bg-warning/10',
      borderGlow: 'hover:border-warning/30',
      shadowGlow: 'hover:shadow-[0_0_30px_rgba(251,191,36,0.15)]',
    },
    {
      labelKey: 'stats.done',
      value: s.by_status.done ?? 0,
      icon: CheckCircle2,
      color: 'text-success',
      bg: 'bg-success/10',
      borderGlow: 'hover:border-success/30',
      shadowGlow: 'hover:shadow-[0_0_30px_rgba(52,211,153,0.15)]',
    },
    {
      labelKey: 'stats.overdue',
      value: s.overdue,
      icon: AlertTriangle,
      color: 'text-danger',
      bg: 'bg-danger/10',
      borderGlow: 'hover:border-danger/30',
      shadowGlow: 'hover:shadow-[0_0_30px_rgba(248,113,113,0.15)]',
    },
  ]
})
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-text-primary mb-1">{{ t('stats.title') }}</h1>
    <p class="text-sm text-text-secondary mb-8">{{ t('stats.subtitle') }}</p>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
      <div
        v-for="card in cards"
        :key="card.labelKey"
        :class="[
          'glow-card rounded-2xl border border-white/[0.08] bg-white/[0.03] p-6 transition-all duration-300 hover:scale-[1.03] hover:-translate-y-1 cursor-default',
          card.borderGlow,
          card.shadowGlow,
        ]"
      >
        <div class="flex items-center justify-between mb-4">
          <span class="text-sm font-medium text-text-secondary">{{ t(card.labelKey) }}</span>
          <div :class="['flex h-10 w-10 items-center justify-center rounded-xl', card.bg]">
            <component :is="card.icon" :class="['w-5 h-5', card.color]" />
          </div>
        </div>
        <p class="text-4xl font-bold tracking-tight text-text-primary animate-count-up">
          {{ card.value }}
        </p>
      </div>
    </div>
  </div>
</template>
