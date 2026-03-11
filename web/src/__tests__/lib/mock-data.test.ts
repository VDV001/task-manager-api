import { describe, it, expect } from 'vitest'
import { generateMockTasks } from '@/lib/mock-data'

describe('generateMockTasks', () => {
  it('generates correct number of tasks', () => {
    const tasks = generateMockTasks(100)
    expect(tasks).toHaveLength(100)
  })

  it('generates default 10000 tasks', () => {
    const tasks = generateMockTasks()
    expect(tasks).toHaveLength(10000)
  })

  it('generates tasks with required fields', () => {
    const tasks = generateMockTasks(5)
    for (const task of tasks) {
      expect(task.id).toBeTruthy()
      expect(task.title).toBeTruthy()
      expect(task.status).toMatch(/^(new|in_progress|done)$/)
      expect(task.created_at).toBeTruthy()
      expect(task.updated_at).toBeTruthy()
      expect(task.author_id).toBeTruthy()
    }
  })

  it('generates unique IDs', () => {
    const tasks = generateMockTasks(50)
    const ids = new Set(tasks.map((t) => t.id))
    expect(ids.size).toBe(50)
  })

  it('generates titles with sequential numbers', () => {
    const tasks = generateMockTasks(5)
    expect(tasks[0]!.title).toContain('#1')
    expect(tasks[4]!.title).toContain('#5')
  })

  it('generates some tasks with deadlines and some without', () => {
    const tasks = generateMockTasks(100)
    const withDeadline = tasks.filter((t) => t.deadline)
    const withoutDeadline = tasks.filter((t) => !t.deadline)
    // With 70% chance of deadline, both groups should have tasks
    expect(withDeadline.length).toBeGreaterThan(0)
    expect(withoutDeadline.length).toBeGreaterThan(0)
  })

  it('generates valid ISO date strings', () => {
    const tasks = generateMockTasks(10)
    for (const task of tasks) {
      expect(new Date(task.created_at).toISOString()).toBe(task.created_at)
      if (task.deadline) {
        expect(new Date(task.deadline).toISOString()).toBe(task.deadline)
      }
    }
  })
})
