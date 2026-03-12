export interface Stock {
  id: number
  ticker: string
  company: string
  brokerage: string
  action: string
  rating_from: string
  rating_to: string
  target_from: number
  target_to: number
  created_at: string
}

export interface StockListResponse {
  items: Stock[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

export interface Recommendation {
  stock: Stock
  score: number
  reason: string
}

export interface RecommendationResponse {
  recommendations: Recommendation[]
  total: number
}

export interface StockFilters {
  search: string
  sortBy: string
  sortOrder: 'asc' | 'desc'
  page: number
  pageSize: number
}
