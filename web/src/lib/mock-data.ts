import type { Task, TaskStatus } from '@/types/task'

const titles = [
  'Fix login page bug',
  'Deploy to production',
  'Update API docs',
  'Refactor auth module',
  'Add unit tests',
  'Setup CI/CD pipeline',
  'Optimize database queries',
  'Design landing page',
  'Implement search feature',
  'Add dark mode support',
  'Fix memory leak',
  'Update dependencies',
  'Write API integration tests',
  'Create user dashboard',
  'Add email notifications',
  'Setup monitoring alerts',
  'Migrate to TypeScript',
  'Add rate limiting',
  'Fix CORS issues',
  'Implement file upload',
  'Add caching layer',
  'Create admin panel',
  'Setup logging',
  'Add pagination support',
  'Fix responsive layout',
  'Implement webhooks',
  'Add OAuth2 support',
  'Create backup script',
  'Setup load balancer',
  'Add data export feature',
]

const descriptions = [
  'This needs to be done ASAP',
  'Low priority but important',
  'Blocking other tasks',
  'Nice to have feature',
  'Critical for release',
  'Follow up from code review',
  'Customer reported issue',
  'Technical debt cleanup',
  'Performance improvement',
  '',
  '',
  '', // some tasks without description
]

const statuses: TaskStatus[] = ['new', 'in_progress', 'done']

function randomDate(start: Date, end: Date): Date {
  return new Date(start.getTime() + Math.random() * (end.getTime() - start.getTime()))
}

function randomId(): string {
  return crypto.randomUUID()
}

export function generateMockTasks(count: number = 10000): Task[] {
  const now = new Date()
  const sixMonthsAgo = new Date(now.getFullYear(), now.getMonth() - 6, 1)
  const oneMonthAhead = new Date(now.getFullYear(), now.getMonth() + 1, now.getDate())

  return Array.from({ length: count }, (_, i) => {
    const createdAt = randomDate(sixMonthsAgo, now)
    const status = statuses[Math.floor(Math.random() * statuses.length)]!
    const hasDeadline = Math.random() > 0.3
    const deadline = hasDeadline ? randomDate(sixMonthsAgo, oneMonthAhead) : undefined

    return {
      id: randomId(),
      title: `${titles[i % titles.length]} #${i + 1}`,
      description: descriptions[Math.floor(Math.random() * descriptions.length)] ?? '',
      status,
      deadline: deadline?.toISOString(),
      created_at: createdAt.toISOString(),
      updated_at: new Date(createdAt.getTime() + Math.random() * 86400000 * 7).toISOString(),
      author_id: randomId(),
    } satisfies Task
  })
}
