import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Stock, StockFilters, Recommendation } from '@/types/stock'
import { fetchStocks, fetchRecommendations, syncStocks, fetchStats } from '@/services/api'

export const useStockStore = defineStore('stocks', () => {
  // State
  const stocks = ref<Stock[]>([])
  const recommendations = ref<Recommendation[]>([])
  const totalStocks = ref(0)
  const totalPages = ref(0)
  const currentPage = ref(1)
  const pageSize = ref(20)
  const isLoading = ref(false)
  const isSyncing = ref(false)
  const error = ref<string | null>(null)
  const statsCount = ref(0)

  const filters = ref<StockFilters>({
    search: '',
    sortBy: 'id',
    sortOrder: 'asc',
    page: 1,
    pageSize: 20,
  })

  // Actions
  async function loadStocks() {
    isLoading.value = true
    error.value = null
    try {
      const result = await fetchStocks({
        search: filters.value.search,
        sort_by: filters.value.sortBy,
        sort_order: filters.value.sortOrder,
        page: filters.value.page,
        page_size: filters.value.pageSize,
      })
      stocks.value = result.items || []
      totalStocks.value = result.total
      totalPages.value = result.total_pages
      currentPage.value = result.page
      pageSize.value = result.page_size
    } catch (e: any) {
      error.value = e.message || 'Failed to load stocks'
      console.error('Error loading stocks:', e)
    } finally {
      isLoading.value = false
    }
  }

  async function loadRecommendations(top: number = 10) {
    isLoading.value = true
    error.value = null
    try {
      const result = await fetchRecommendations(top)
      recommendations.value = result.recommendations || []
    } catch (e: any) {
      error.value = e.message || 'Failed to load recommendations'
      console.error('Error loading recommendations:', e)
    } finally {
      isLoading.value = false
    }
  }

  async function triggerSync() {
    isSyncing.value = true
    error.value = null
    try {
      const result = await syncStocks()
      await loadStocks()
      return result
    } catch (e: any) {
      error.value = e.message || 'Sync failed'
      console.error('Error syncing:', e)
    } finally {
      isSyncing.value = false
    }
  }

  async function loadStats() {
    try {
      const result = await fetchStats()
      statsCount.value = result.total_stocks
    } catch (e: any) {
      console.error('Error loading stats:', e)
    }
  }

  function setSearch(search: string) {
    filters.value.search = search
    filters.value.page = 1
    loadStocks()
  }

  function setSort(column: string) {
    if (filters.value.sortBy === column) {
      filters.value.sortOrder = filters.value.sortOrder === 'asc' ? 'desc' : 'asc'
    } else {
      filters.value.sortBy = column
      filters.value.sortOrder = 'asc'
    }
    loadStocks()
  }

  function setPage(page: number) {
    filters.value.page = page
    loadStocks()
  }

  return {
    stocks,
    recommendations,
    totalStocks,
    totalPages,
    currentPage,
    pageSize,
    isLoading,
    isSyncing,
    error,
    filters,
    statsCount,
    loadStocks,
    loadRecommendations,
    triggerSync,
    loadStats,
    setSearch,
    setSort,
    setPage,
  }
})
