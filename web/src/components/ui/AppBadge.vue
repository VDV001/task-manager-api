<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { cn } from '@/lib/utils'
import type { TaskStatus } from '@/types/task'

const { t } = useI18n()

const props = defineProps<{
  status: TaskStatus
}>()

const config = computed(() => {
  switch (props.status) {
    case 'new':
      return { labelKey: 'task.new', class: 'bg-info/10 text-info border-info/20' }
    case 'in_progress':
      return { labelKey: 'task.inProgress', class: 'bg-warning/10 text-warning border-warning/20' }
    case 'done':
      return { labelKey: 'task.done', class: 'bg-success/10 text-success border-success/20' }
    default:
      return { labelKey: '', class: '' }
  }
})
</script>

<template>
  <span
    :class="
      cn(
        'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border',
        config.class,
      )
    "
  >
    {{ config.labelKey ? t(config.labelKey) : status }}
  </span>
</template>
