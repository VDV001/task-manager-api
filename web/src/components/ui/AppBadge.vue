<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/lib/utils'
import type { TaskStatus } from '@/types/task'

const props = defineProps<{
  status: TaskStatus
}>()

const config = computed(() => {
  switch (props.status) {
    case 'new':
      return { label: 'New', class: 'bg-info/10 text-info border-info/20' }
    case 'in_progress':
      return { label: 'In Progress', class: 'bg-warning/10 text-warning border-warning/20' }
    case 'done':
      return { label: 'Done', class: 'bg-success/10 text-success border-success/20' }
  }
})
</script>

<template>
  <span
    :class="cn(
      'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border',
      config.class,
    )"
  >
    {{ config.label }}
  </span>
</template>
