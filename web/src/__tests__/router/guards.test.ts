import { describe, it, expect, beforeEach, vi } from 'vitest'
import { createRouter, createWebHistory } from 'vue-router'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '@/stores/auth'

// Mock API client
vi.mock('@/api/client', () => ({
  setupApiClient: vi.fn(),
}))

vi.mock('@/api/auth', () => ({
  authApi: {
    login: vi.fn(),
    register: vi.fn(),
    refresh: vi.fn(),
  },
}))

function createTestRouter() {
  const router = createRouter({
    history: createWebHistory(),
    routes: [
      {
        path: '/login',
        name: 'login',
        component: { template: '<div>Login</div>' },
        meta: { guest: true },
      },
      {
        path: '/register',
        name: 'register',
        component: { template: '<div>Register</div>' },
        meta: { guest: true },
      },
      {
        path: '/',
        component: { template: '<router-view />' },
        meta: { requiresAuth: true },
        children: [
          { path: '', name: 'dashboard', component: { template: '<div>Dashboard</div>' } },
          { path: 'stats', name: 'stats', component: { template: '<div>Stats</div>' } },
        ],
      },
    ],
  })

  router.beforeEach((to) => {
    const auth = useAuthStore()
    if (to.meta.requiresAuth && !auth.isAuthenticated) {
      return { name: 'login' }
    }
    if (to.meta.guest && auth.isAuthenticated) {
      return { name: 'dashboard' }
    }
  })

  return router
}

describe('Router Guards', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it('redirects unauthenticated user from dashboard to login', async () => {
    const router = createTestRouter()
    await router.push('/')
    await router.isReady()

    expect(router.currentRoute.value.name).toBe('login')
  })

  it('redirects unauthenticated user from stats to login', async () => {
    const router = createTestRouter()
    await router.push('/stats')
    await router.isReady()

    expect(router.currentRoute.value.name).toBe('login')
  })

  it('allows authenticated user to access dashboard', async () => {
    localStorage.setItem('access_token', 'test-token')
    localStorage.setItem('refresh_token', 'test-refresh')

    const router = createTestRouter()
    await router.push('/')
    await router.isReady()

    expect(router.currentRoute.value.name).toBe('dashboard')
  })

  it('redirects authenticated user from login to dashboard', async () => {
    localStorage.setItem('access_token', 'test-token')
    localStorage.setItem('refresh_token', 'test-refresh')

    const router = createTestRouter()
    await router.push('/login')
    await router.isReady()

    expect(router.currentRoute.value.name).toBe('dashboard')
  })

  it('allows unauthenticated user to access login', async () => {
    const router = createTestRouter()
    await router.push('/login')
    await router.isReady()

    expect(router.currentRoute.value.name).toBe('login')
  })
})
