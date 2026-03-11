<script setup lang="ts">
import { computed } from 'vue'
import { Doughnut } from 'vue-chartjs'
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js'
import type { TaskStats } from '@/types/task'

ChartJS.register(ArcElement, Tooltip, Legend)

const props = defineProps<{
  stats: TaskStats | null
}>()

const chartData = computed(() => {
  const s = props.stats
  if (!s) return null

  return {
    labels: ['New', 'In Progress', 'Done'],
    datasets: [
      {
        data: [s.by_status.new ?? 0, s.by_status.in_progress ?? 0, s.by_status.done ?? 0],
        backgroundColor: [
          'rgba(56, 189, 248, 0.8)',
          'rgba(251, 191, 36, 0.8)',
          'rgba(52, 211, 153, 0.8)',
        ],
        borderColor: ['#38bdf8', '#fbbf24', '#34d399'],
        borderWidth: 2,
        hoverOffset: 8,
      },
    ],
  }
})

const options = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'bottom' as const,
      labels: {
        color: '#8888a0',
        padding: 16,
        usePointStyle: true,
        pointStyleWidth: 10,
      },
    },
    tooltip: {
      backgroundColor: '#16161f',
      borderColor: 'rgba(255,255,255,0.08)',
      borderWidth: 1,
      titleColor: '#f0f0f5',
      bodyColor: '#8888a0',
    },
  },
  cutout: '60%',
}
</script>

<template>
  <div class="h-64">
    <Doughnut v-if="chartData" :data="chartData" :options="options" />
    <div v-else class="flex items-center justify-center h-full text-text-muted text-sm">
      No data to display
    </div>
  </div>
</template>
