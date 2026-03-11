<script setup lang="ts">
import { computed } from 'vue'
import { Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Filler,
} from 'chart.js'
import type { Task } from '@/types/task'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Filler)

const props = defineProps<{
  tasks: Task[]
}>()

const chartData = computed(() => {
  // Group tasks by month
  const months: Record<string, number> = {}
  props.tasks.forEach((task) => {
    const date = new Date(task.created_at)
    const key = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`
    months[key] = (months[key] || 0) + 1
  })

  const sorted = Object.entries(months).sort(([a], [b]) => a.localeCompare(b))
  const labels = sorted.map(([k]) => {
    const [y, m] = k.split('-')
    return new Date(Number(y), Number(m) - 1).toLocaleDateString('en', {
      month: 'short',
      year: '2-digit',
    })
  })
  const data = sorted.map(([, v]) => v)

  return {
    labels,
    datasets: [
      {
        label: 'Tasks Created',
        data,
        borderColor: '#6366f1',
        backgroundColor: 'rgba(99, 102, 241, 0.1)',
        fill: true,
        tension: 0.4,
        pointBackgroundColor: '#6366f1',
        pointBorderColor: '#6366f1',
        pointRadius: 4,
        pointHoverRadius: 6,
      },
    ],
  }
})

const options = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    title: { display: false },
    tooltip: {
      backgroundColor: '#16161f',
      borderColor: 'rgba(255,255,255,0.08)',
      borderWidth: 1,
      titleColor: '#f0f0f5',
      bodyColor: '#8888a0',
    },
  },
  scales: {
    x: {
      grid: { color: 'rgba(255,255,255,0.04)' },
      ticks: { color: '#55556a' },
    },
    y: {
      grid: { color: 'rgba(255,255,255,0.04)' },
      ticks: { color: '#55556a' },
      beginAtZero: true,
    },
  },
}
</script>

<template>
  <div class="h-64">
    <Line v-if="chartData.labels.length" :data="chartData" :options="options" />
    <div v-else class="flex items-center justify-center h-full text-text-muted text-sm">
      No data to display
    </div>
  </div>
</template>
