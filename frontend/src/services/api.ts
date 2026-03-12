import axios from 'axios'
import type { StockListResponse, RecommendationResponse, Stock } from '@/types/stock'

const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  headers: {
    'Content-Type': 'application/json',
  },
})

export async function fetchStocks(params: {
  search?: string
  sort_by?: string
  sort_order?: string
  page?: number
  page_size?: number
}): Promise<StockListResponse> {
  const { data } = await api.get<StockListResponse>('/stocks', { params })
  return data
}

export async function fetchStockById(id: number): Promise<Stock> {
  const { data } = await api.get<Stock>(`/stocks/${id}`)
  return data
}

export async function fetchRecommendations(top: number = 10): Promise<RecommendationResponse> {
  const { data } = await api.get<RecommendationResponse>('/recommendations', {
    params: { top },
  })
  return data
}

export async function syncStocks(): Promise<{ message: string; items_synced: number }> {
  const { data } = await api.post('/sync')
  return data
}

export async function fetchStats(): Promise<{ total_stocks: number }> {
  const { data } = await api.get('/stats')
  return data
}
