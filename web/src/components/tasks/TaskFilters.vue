<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { Search } from 'lucide-vue-next'
import type { TaskFilter, TaskStatus } from '@/types/task'

const { t } = useI18n()

const props = defineProps<{
  filter: TaskFilter
}>()

const emit = defineEmits<{
  update: [filter: Partial<TaskFilter>]
}>()

const search = ref(props.filter.search ?? '')
let searchTimeout: ReturnType<typeof setTimeout>

watch(search, (val) => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    emit('update', { search: val || undefined })
  }, 300)
})

const statuses: { value: TaskStatus | undefined; labelKey: string }[] = [
  { value: undefined, labelKey: 'dashboard.allStatuses' },
  { value: 'new', labelKey: 'task.new' },
  { value: 'in_progress', labelKey: 'task.inProgress' },
  { value: 'done', labelKey: 'task.done' },
]
</script>

<template>
  <div class="flex flex-wrap items-center gap-3">
    <!-- Search -->
    <div class="relative flex-1 min-w-[200px]">
      <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-text-muted" />
      <input
        v-model="search"
        type="text"
        :placeholder="t('dashboard.searchTasks')"
        class="w-full pl-10 pr-4 py-2.5 rounded-xl bg-white/5 border border-border text-sm text-text-primary placeholder-text-muted focus:outline-none focus:ring-2 focus:ring-accent/50 focus:border-accent transition-all"
      />
    </div>

    <!-- Status filter -->
    <div class="flex gap-1 bg-white/5 rounded-xl p-1 border border-border">
      <button
        v-for="s in statuses"
        :key="s.labelKey"
        :class="[
          'px-3 py-1.5 rounded-lg text-xs font-medium transition-all cursor-pointer',
          filter.status === s.value
            ? 'bg-accent/15 text-accent'
            : 'text-text-secondary hover:text-text-primary hover:bg-white/5',
        ]"
        @click="$emit('update', { status: s.value })"
      >
        {{ t(s.labelKey) }}
      </button>
    </div>

    <!-- Overdue toggle -->
    <button
      :class="[
        'px-3 py-2 rounded-xl text-xs font-medium border transition-all cursor-pointer',
        filter.overdue
          ? 'bg-danger/10 text-danger border-danger/20'
          : 'bg-white/5 text-text-secondary border-border hover:text-text-primary',
      ]"
      @click="$emit('update', { overdue: !filter.overdue || undefined })"
    >
      {{ t('dashboard.overdue') }}
    </button>
  </div>
</template>
