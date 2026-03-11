export interface ApiResponse<T> {
  data: T
  meta?: PaginationMeta
}

export interface PaginationMeta {
  page: number
  limit: number
  total: number
}

export interface ApiError {
  error: {
    code: string
    message: string
    details?: FieldError[]
  }
}

export interface FieldError {
  field: string
  message: string
}
