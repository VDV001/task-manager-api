import { api } from './client'
import type {
  Task,
  CreateTaskRequest,
  UpdateTaskRequest,
  TaskStats,
  TaskFilter,
} from '@/types/task'
import type { ApiResponse } from '@/types/api'

export const tasksApi = {
  list(filter?: TaskFilter) {
    const query: Record<string, string> = {}
    if (filter) {
      Object.entries(filter).forEach(([key, value]) => {
        if (value !== undefined && value !== '') {
          query[key] = String(value)
        }
      })
    }
    return api<ApiResponse<Task[]>>('/tasks', { query })
  },

  getById(id: string) {
    return api<ApiResponse<Task>>(`/tasks/${id}`)
  },

  create(data: CreateTaskRequest) {
    return api<ApiResponse<Task>>('/tasks', {
      method: 'POST',
      body: data,
    })
  },

  update(id: string, data: UpdateTaskRequest) {
    return api<ApiResponse<Task>>(`/tasks/${id}`, {
      method: 'PATCH',
      body: data,
    })
  },

  delete(id: string) {
    return api<void>(`/tasks/${id}`, { method: 'DELETE' })
  },

  stats() {
    return api<ApiResponse<TaskStats>>('/tasks/stats')
  },
}
