import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useTaskStore } from '@/stores/tasks'

vi.mock('@/api/tasks', () => ({
  tasksApi: {
    list: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    delete: vi.fn(),
    stats: vi.fn(),
  },
}))

import { tasksApi } from '@/api/tasks'

const mockTask = {
  id: '1',
  title: 'Test task',
  description: 'Desc',
  status: 'new' as const,
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
  author_id: 'user-1',
}

const mockStats = {
  total: 5,
  by_status: { new: 2, in_progress: 2, done: 1 },
  overdue: 1,
}

describe('Task Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('has correct initial state', () => {
    const store = useTaskStore()
    expect(store.tasks).toEqual([])
    expect(store.stats).toBeNull()
    expect(store.loading).toBe(false)
    expect(store.filter.page).toBe(1)
    expect(store.filter.limit).toBe(12)
    expect(store.totalPages).toBe(1)
  })

  it('fetchTasks loads tasks and meta', async () => {
    vi.mocked(tasksApi.list).mockResolvedValue({
      data: [mockTask],
      meta: { page: 1, limit: 12, total: 25 },
    })

    const store = useTaskStore()
    await store.fetchTasks()

    expect(store.tasks).toHaveLength(1)
    expect(store.tasks[0]!.title).toBe('Test task')
    expect(store.meta!.total).toBe(25)
    expect(store.loading).toBe(false)
  })

  it('totalPages computes from meta', async () => {
    vi.mocked(tasksApi.list).mockResolvedValue({
      data: [mockTask],
      meta: { page: 1, limit: 12, total: 25 },
    })

    const store = useTaskStore()
    await store.fetchTasks()

    expect(store.totalPages).toBe(3) // ceil(25/12) = 3
  })

  it('fetchTasks sets loading during fetch', async () => {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    let resolve: (val: any) => void
    vi.mocked(tasksApi.list).mockReturnValue(
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      new Promise<any>((r) => {
        resolve = r
      }),
    )

    const store = useTaskStore()
    const promise = store.fetchTasks()

    expect(store.loading).toBe(true)

    resolve!({ data: [], meta: undefined })
    await promise

    expect(store.loading).toBe(false)
  })

  it('fetchStats loads stats', async () => {
    vi.mocked(tasksApi.stats).mockResolvedValue({ data: mockStats })

    const store = useTaskStore()
    await store.fetchStats()

    expect(store.stats).toEqual(mockStats)
  })

  it('createTask adds task to list and increments total', async () => {
    vi.mocked(tasksApi.list).mockResolvedValue({
      data: [],
      meta: { page: 1, limit: 12, total: 0 },
    })
    vi.mocked(tasksApi.create).mockResolvedValue({ data: mockTask })
    vi.mocked(tasksApi.stats).mockResolvedValue({ data: mockStats })

    const store = useTaskStore()
    await store.fetchTasks()

    const task = await store.createTask({ title: 'Test task' })

    expect(task.title).toBe('Test task')
    expect(store.tasks).toHaveLength(1)
    expect(store.tasks[0]).toStrictEqual(task)
    expect(store.meta!.total).toBe(1)
  })

  it('updateTask replaces task in list', async () => {
    vi.mocked(tasksApi.list).mockResolvedValue({
      data: [mockTask],
      meta: { page: 1, limit: 12, total: 1 },
    })
    const updatedTask = { ...mockTask, title: 'Updated' }
    vi.mocked(tasksApi.update).mockResolvedValue({ data: updatedTask })
    vi.mocked(tasksApi.stats).mockResolvedValue({ data: mockStats })

    const store = useTaskStore()
    await store.fetchTasks()
    await store.updateTask('1', { title: 'Updated' })

    expect(store.tasks[0]!.title).toBe('Updated')
  })

  it('deleteTask removes task and decrements total', async () => {
    vi.mocked(tasksApi.list).mockResolvedValue({
      data: [mockTask],
      meta: { page: 1, limit: 12, total: 1 },
    })
    vi.mocked(tasksApi.delete).mockResolvedValue(undefined as never)
    vi.mocked(tasksApi.stats).mockResolvedValue({ data: mockStats })

    const store = useTaskStore()
    await store.fetchTasks()
    await store.deleteTask('1')

    expect(store.tasks).toHaveLength(0)
    expect(store.meta!.total).toBe(0)
  })

  it('setFilter updates filter and refetches', async () => {
    vi.mocked(tasksApi.list).mockResolvedValue({ data: [] })

    const store = useTaskStore()
    store.setFilter({ status: 'done' })

    expect(store.filter.status).toBe('done')
    expect(store.filter.page).toBe(1) // resets page
    expect(tasksApi.list).toHaveBeenCalled()
  })

  it('setPage updates page and refetches', async () => {
    vi.mocked(tasksApi.list).mockResolvedValue({ data: [] })

    const store = useTaskStore()
    store.setPage(3)

    expect(store.filter.page).toBe(3)
    expect(tasksApi.list).toHaveBeenCalled()
  })
})
