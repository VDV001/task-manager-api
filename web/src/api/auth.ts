import { api } from './client'
import type { LoginRequest, RegisterRequest, RefreshRequest, TokenResponse } from '@/types/auth'
import type { ApiResponse } from '@/types/api'

export const authApi = {
  login(data: LoginRequest) {
    return api<ApiResponse<TokenResponse>>('/auth/login', {
      method: 'POST',
      body: data,
    })
  },

  register(data: RegisterRequest) {
    return api<ApiResponse<TokenResponse>>('/auth/register', {
      method: 'POST',
      body: data,
    })
  },

  refresh(data: RefreshRequest) {
    return api<ApiResponse<TokenResponse>>('/auth/refresh', {
      method: 'POST',
      body: data,
    })
  },
}
