export type TaskStatus = 'new' | 'in_progress' | 'done'

export interface Task {
  id: string
  title: string
  description: string
  status: TaskStatus
  deadline?: string
  created_at: string
  updated_at: string
  author_id: string
}

export interface CreateTaskRequest {
  title: string
  description?: string
  deadline?: string
}

export interface UpdateTaskRequest {
  title?: string
  description?: string
  status?: TaskStatus
  deadline?: string
}

export interface TaskStats {
  total: number
  by_status: Record<string, number>
  overdue: number
}

export interface TaskFilter {
  status?: TaskStatus
  search?: string
  overdue?: boolean
  deadline_before?: string
  deadline_after?: string
  sort_by?: 'created_at' | 'deadline' | 'status' | 'title'
  order?: 'asc' | 'desc'
  page?: number
  limit?: number
}
