<script setup lang="ts">
import { computed } from 'vue'
import { formatDistanceToNow, isPast, parseISO } from 'date-fns'
import { Calendar, Pencil, Trash2 } from 'lucide-vue-next'
import AppCard from '@/components/ui/AppCard.vue'
import AppBadge from '@/components/ui/AppBadge.vue'
import type { Task } from '@/types/task'

const props = defineProps<{
  task: Task
}>()

defineEmits<{
  edit: [task: Task]
  delete: [task: Task]
}>()

const isOverdue = computed(() => {
  if (!props.task.deadline || props.task.status === 'done') return false
  return isPast(parseISO(props.task.deadline))
})

const deadlineText = computed(() => {
  if (!props.task.deadline) return null
  return formatDistanceToNow(parseISO(props.task.deadline), { addSuffix: true })
})
</script>

<template>
  <AppCard hoverable class="group">
    <div class="flex items-start justify-between mb-3">
      <AppBadge :status="task.status" />
      <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
        <button
          class="p-1.5 rounded-lg hover:bg-white/10 text-text-muted hover:text-text-primary transition-colors cursor-pointer"
          @click.stop="$emit('edit', task)"
        >
          <Pencil class="w-4 h-4" />
        </button>
        <button
          class="p-1.5 rounded-lg hover:bg-danger/10 text-text-muted hover:text-danger transition-colors cursor-pointer"
          @click.stop="$emit('delete', task)"
        >
          <Trash2 class="w-4 h-4" />
        </button>
      </div>
    </div>

    <h3 class="text-base font-semibold text-text-primary mb-1 line-clamp-1">
      {{ task.title }}
    </h3>
    <p v-if="task.description" class="text-sm text-text-secondary line-clamp-2 mb-3">
      {{ task.description }}
    </p>

    <div v-if="deadlineText" class="flex items-center gap-1.5 text-xs" :class="isOverdue ? 'text-danger' : 'text-text-muted'">
      <Calendar class="w-3.5 h-3.5" />
      <span>{{ deadlineText }}</span>
    </div>
  </AppCard>
</template>
