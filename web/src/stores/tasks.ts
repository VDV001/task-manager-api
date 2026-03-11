import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { tasksApi } from '@/api/tasks'
import type { Task, TaskFilter, TaskStats } from '@/types/task'
import type { PaginationMeta } from '@/types/api'

export const useTaskStore = defineStore('tasks', () => {
  const tasks = ref<Task[]>([])
  const stats = ref<TaskStats | null>(null)
  const meta = ref<PaginationMeta | null>(null)
  const loading = ref(false)
  const filter = ref<TaskFilter>({
    page: 1,
    limit: 12,
    sort_by: 'created_at',
    order: 'desc',
  })

  const totalPages = computed(() => {
    if (!meta.value) return 1
    return Math.ceil(meta.value.total / meta.value.limit)
  })

  async function fetchTasks() {
    loading.value = true
    try {
      const res = await tasksApi.list(filter.value)
      tasks.value = res.data
      meta.value = res.meta ?? null
    } finally {
      loading.value = false
    }
  }

  async function fetchStats() {
    const res = await tasksApi.stats()
    stats.value = res.data
  }

  async function createTask(data: Parameters<typeof tasksApi.create>[0]) {
    const res = await tasksApi.create(data)
    tasks.value.unshift(res.data)
    if (meta.value) meta.value.total++
    await fetchStats()
    return res.data
  }

  async function updateTask(id: string, data: Parameters<typeof tasksApi.update>[1]) {
    const res = await tasksApi.update(id, data)
    const idx = tasks.value.findIndex((t) => t.id === id)
    if (idx !== -1) tasks.value[idx] = res.data
    await fetchStats()
    return res.data
  }

  async function deleteTask(id: string) {
    await tasksApi.delete(id)
    tasks.value = tasks.value.filter((t) => t.id !== id)
    if (meta.value) meta.value.total--
    await fetchStats()
  }

  function setFilter(newFilter: Partial<TaskFilter>) {
    filter.value = { ...filter.value, ...newFilter, page: newFilter.page ?? 1 }
    fetchTasks()
  }

  function setPage(page: number) {
    filter.value.page = page
    fetchTasks()
  }

  return {
    tasks, stats, meta, loading, filter, totalPages,
    fetchTasks, fetchStats, createTask, updateTask, deleteTask,
    setFilter, setPage,
  }
})
