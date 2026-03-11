import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { authApi } from '@/api/auth'
import { setupApiClient } from '@/api/client'
import type { LoginRequest, RegisterRequest } from '@/types/auth'

export const useAuthStore = defineStore('auth', () => {
  const accessToken = ref<string | null>(localStorage.getItem('access_token'))
  const refreshToken = ref<string | null>(localStorage.getItem('refresh_token'))

  const isAuthenticated = computed(() => !!accessToken.value)

  function setTokens(access: string, refresh: string) {
    accessToken.value = access
    refreshToken.value = refresh
    localStorage.setItem('access_token', access)
    localStorage.setItem('refresh_token', refresh)
  }

  function clearTokens() {
    accessToken.value = null
    refreshToken.value = null
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
  }

  async function login(data: LoginRequest) {
    const res = await authApi.login(data)
    setTokens(res.data.access_token, res.data.refresh_token)
  }

  async function register(data: RegisterRequest) {
    const res = await authApi.register(data)
    setTokens(res.data.access_token, res.data.refresh_token)
  }

  async function refresh() {
    if (!refreshToken.value) {
      clearTokens()
      return false
    }
    try {
      const res = await authApi.refresh({ refresh_token: refreshToken.value })
      setTokens(res.data.access_token, res.data.refresh_token)
      return true
    } catch {
      clearTokens()
      return false
    }
  }

  function logout() {
    clearTokens()
  }

  // Setup API client interceptors
  setupApiClient(
    () => accessToken.value,
    async () => {
      const refreshed = await refresh()
      if (!refreshed) {
        clearTokens()
        window.location.href = '/login'
      }
    },
  )

  return {
    accessToken,
    refreshToken,
    isAuthenticated,
    login,
    register,
    refresh,
    logout,
  }
})
