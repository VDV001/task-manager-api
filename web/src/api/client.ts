import { ofetch } from 'ofetch'

let getAccessToken: (() => string | null) | null = null
let onUnauthorized: (() => void) | null = null

export function setupApiClient(tokenGetter: () => string | null, unauthorizedHandler: () => void) {
  getAccessToken = tokenGetter
  onUnauthorized = unauthorizedHandler
}

export const api = ofetch.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  onRequest({ options }) {
    const token = getAccessToken?.()
    if (token) {
      const headers = new Headers(options.headers)
      headers.set('Authorization', `Bearer ${token}`)
      options.headers = headers
    }
  },
  onResponseError({ response }) {
    if (response.status === 401) {
      onUnauthorized?.()
    }
  },
})
