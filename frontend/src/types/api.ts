export type ApiResponse<T> = {
  success: boolean
  message: string
  data: T
}

export type TokenPair = {
  access_token: string
  refresh_token: string
}

export type PageMeta = {
  page: number
  limit: number
  total: number
  total_pages: number
}

export type Paginated<T> = {
  items: T[]
  meta: PageMeta
}
