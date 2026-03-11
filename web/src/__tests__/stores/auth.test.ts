import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '@/stores/auth'

// Mock API client setup and auth API
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

import { authApi } from '@/api/auth'

describe('Auth Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    vi.clearAllMocks()
  })

  it('starts unauthenticated when no tokens in localStorage', () => {
    const auth = useAuthStore()
    expect(auth.isAuthenticated).toBe(false)
    expect(auth.accessToken).toBeNull()
    expect(auth.refreshToken).toBeNull()
  })

  it('starts authenticated when tokens exist in localStorage', () => {
    localStorage.setItem('access_token', 'test-access')
    localStorage.setItem('refresh_token', 'test-refresh')

    const auth = useAuthStore()
    expect(auth.isAuthenticated).toBe(true)
    expect(auth.accessToken).toBe('test-access')
  })

  it('login stores tokens and sets authenticated', async () => {
    vi.mocked(authApi.login).mockResolvedValue({
      data: { access_token: 'new-access', refresh_token: 'new-refresh' },
    })

    const auth = useAuthStore()
    await auth.login({ email: 'test@test.com', password: '123456' })

    expect(auth.isAuthenticated).toBe(true)
    expect(auth.accessToken).toBe('new-access')
    expect(auth.refreshToken).toBe('new-refresh')
    expect(localStorage.setItem).toHaveBeenCalledWith('access_token', 'new-access')
    expect(localStorage.setItem).toHaveBeenCalledWith('refresh_token', 'new-refresh')
  })

  it('register stores tokens', async () => {
    vi.mocked(authApi.register).mockResolvedValue({
      data: { access_token: 'reg-access', refresh_token: 'reg-refresh' },
    })

    const auth = useAuthStore()
    await auth.register({ name: 'Test', email: 'test@test.com', password: '123456' })

    expect(auth.isAuthenticated).toBe(true)
    expect(auth.accessToken).toBe('reg-access')
  })

  it('logout clears tokens', async () => {
    vi.mocked(authApi.login).mockResolvedValue({
      data: { access_token: 'access', refresh_token: 'refresh' },
    })

    const auth = useAuthStore()
    await auth.login({ email: 'test@test.com', password: '123456' })
    expect(auth.isAuthenticated).toBe(true)

    auth.logout()
    expect(auth.isAuthenticated).toBe(false)
    expect(auth.accessToken).toBeNull()
    expect(localStorage.removeItem).toHaveBeenCalledWith('access_token')
    expect(localStorage.removeItem).toHaveBeenCalledWith('refresh_token')
  })

  it('refresh updates tokens on success', async () => {
    localStorage.setItem('refresh_token', 'old-refresh')

    vi.mocked(authApi.refresh).mockResolvedValue({
      data: { access_token: 'refreshed-access', refresh_token: 'refreshed-refresh' },
    })

    const auth = useAuthStore()
    const result = await auth.refresh()

    expect(result).toBe(true)
    expect(auth.accessToken).toBe('refreshed-access')
  })

  it('refresh clears tokens on failure', async () => {
    localStorage.setItem('access_token', 'old-access')
    localStorage.setItem('refresh_token', 'old-refresh')

    vi.mocked(authApi.refresh).mockRejectedValue(new Error('expired'))

    const auth = useAuthStore()
    const result = await auth.refresh()

    expect(result).toBe(false)
    expect(auth.isAuthenticated).toBe(false)
  })

  it('refresh returns false when no refresh token', async () => {
    const auth = useAuthStore()
    const result = await auth.refresh()

    expect(result).toBe(false)
    expect(authApi.refresh).not.toHaveBeenCalled()
  })
})
