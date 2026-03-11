<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useTaskStore } from '@/stores/tasks'
import { ListTodo, Clock, CheckCircle2, AlertTriangle } from 'lucide-vue-next'
import AppCard from '@/components/ui/AppCard.vue'

const taskStore = useTaskStore()

onMounted(() => {
  taskStore.fetchStats()
})

const cards = computed(() => {
  const s = taskStore.stats
  if (!s) return []
  return [
    {
      label: 'Total Tasks',
      value: s.total,
      icon: ListTodo,
      color: 'text-accent',
      bg: 'bg-accent/10',
      glow: 'shadow-accent-glow',
    },
    {
      label: 'New',
      value: s.by_status.new ?? 0,
      icon: Clock,
      color: 'text-info',
      bg: 'bg-info/10',
      glow: 'shadow-[0_0_20px_rgba(56,189,248,0.15)]',
    },
    {
      label: 'In Progress',
      value: s.by_status.in_progress ?? 0,
      icon: AlertTriangle,
      color: 'text-warning',
      bg: 'bg-warning/10',
      glow: 'shadow-[0_0_20px_rgba(251,191,36,0.15)]',
    },
    {
      label: 'Done',
      value: s.by_status.done ?? 0,
      icon: CheckCircle2,
      color: 'text-success',
      bg: 'bg-success/10',
      glow: 'shadow-[0_0_20px_rgba(52,211,153,0.15)]',
    },
    {
      label: 'Overdue',
      value: s.overdue,
      icon: AlertTriangle,
      color: 'text-danger',
      bg: 'bg-danger/10',
      glow: 'shadow-[0_0_20px_rgba(248,113,113,0.15)]',
    },
  ]
})
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold text-text-primary mb-2">Statistics</h1>
    <p class="text-sm text-text-secondary mb-8">Overview of your task progress</p>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
      <AppCard
        v-for="card in cards"
        :key="card.label"
        :class="card.glow"
        hoverable
      >
        <div class="flex items-center justify-between mb-4">
          <span class="text-sm font-medium text-text-secondary">{{ card.label }}</span>
          <div :class="['p-2 rounded-xl', card.bg]">
            <component :is="card.icon" :class="['w-5 h-5', card.color]" />
          </div>
        </div>
        <p class="text-4xl font-bold tracking-tight text-text-primary">{{ card.value }}</p>
      </AppCard>
    </div>
  </div>
</template>
