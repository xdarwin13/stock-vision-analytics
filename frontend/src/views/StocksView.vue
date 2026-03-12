<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useStockStore } from '@/stores/stockStore'

const store = useStockStore()
const searchInput = ref('')
let debounceTimer: ReturnType<typeof setTimeout>
const selectedStock = ref<any>(null)

onMounted(() => {
  store.loadStocks()
})

function handleSearch() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    store.setSearch(searchInput.value)
  }, 300)
}

function getActionBadgeClass(action: string): string {
  const normalized = action.toLowerCase()
  if (normalized.includes('upgraded')) return 'badge-upgraded'
  if (normalized.includes('downgraded')) return 'badge-downgraded'
  if (normalized.includes('maintained') || normalized.includes('reiterated')) return 'badge-maintained'
  if (normalized.includes('initiated')) return 'badge-initiated'
  if (normalized.includes('target raised')) return 'badge-upgraded'
  if (normalized.includes('target lowered')) return 'badge-downgraded'
  return 'badge-default'
}

function formatPrice(value: number): string {
  if (!value) return '—'
  return `$${value.toFixed(2)}`
}

function getSortIcon(column: string): string {
  if (store.filters.sortBy !== column) return ''
  return store.filters.sortOrder === 'asc' ? ' (asc)' : ' (desc)'
}

function getTargetChangeClass(from: number, to: number): string {
  if (!from || !to) return ''
  if (to > from) return 'text-accent-400'
  if (to < from) return 'text-danger-400'
  return ''
}

const columns = [
  { key: 'ticker', label: 'Ticker' },
  { key: 'company', label: 'Company' },
  { key: 'brokerage', label: 'Brokerage' },
  { key: 'action', label: 'Action' },
  { key: 'rating_from', label: 'Rating From' },
  { key: 'rating_to', label: 'Rating To' },
  { key: 'target_from', label: 'Target From' },
  { key: 'target_to', label: 'Target To' },
]
</script>

<template>
  <div>
    <!-- Header -->
    <div class="mb-8 animate-fade-in-up">
      <h1 class="text-3xl font-bold text-white mb-2">Stock Analyst Ratings</h1>
      <p class="text-surface-200/50">Browse, search, and sort analyst ratings and price targets</p>
    </div>

    <!-- Controls -->
    <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-4 mb-6 animate-fade-in-up" style="animation-delay: 0.1s">
      <!-- Search -->
      <div class="relative flex-1">
        <svg class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-surface-200/40" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <input
          v-model="searchInput"
          @input="handleSearch"
          type="text"
          placeholder="Search by ticker, company, brokerage..."
          class="w-full pl-11 pr-4 py-3 rounded-xl bg-surface-800/50 border border-white/10 text-white placeholder-surface-200/30 focus:outline-none focus:border-primary-500/50 focus:ring-2 focus:ring-primary-500/20 transition-all text-sm"
        />
      </div>
      <!-- Sync Button -->
      <button
        @click="store.triggerSync()"
        :disabled="store.isSyncing"
        class="btn-primary whitespace-nowrap"
      >
        <svg v-if="store.isSyncing" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
        </svg>
        <span>{{ store.isSyncing ? 'Syncing...' : 'Sync Data' }}</span>
      </button>
    </div>

    <!-- Results count -->
    <div class="flex items-center justify-between mb-4 text-sm text-surface-200/40">
      <span>{{ store.totalStocks }} results found</span>
      <span>Page {{ store.currentPage }} of {{ store.totalPages }}</span>
    </div>

    <!-- Table -->
    <div class="glass-card overflow-hidden mb-6 animate-fade-in-up" style="animation-delay: 0.2s">
      <!-- Loading -->
      <div v-if="store.isLoading" class="flex items-center justify-center py-20">
        <div class="flex items-center gap-3 text-surface-200/50">
          <svg class="animate-spin w-5 h-5" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
          </svg>
          <span>Loading stocks...</span>
        </div>
      </div>

      <!-- Data Table -->
      <div v-else class="overflow-x-auto">
        <table class="data-table">
          <thead>
            <tr>
              <th
                v-for="col in columns"
                :key="col.key"
                @click="store.setSort(col.key)"
                :class="{ active: store.filters.sortBy === col.key }"
              >
                <span class="flex items-center gap-1">
                  {{ col.label }}
                  <span class="text-[10px]">{{ getSortIcon(col.key) }}</span>
                </span>
              </th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="stock in store.stocks"
              :key="stock.id"
              @click="selectedStock = stock"
              class="cursor-pointer"
            >
              <td class="font-bold text-white">{{ stock.ticker }}</td>
              <td class="text-surface-200/70 max-w-[200px] truncate">{{ stock.company }}</td>
              <td class="text-surface-200/70">{{ stock.brokerage }}</td>
              <td>
                <span :class="['badge', getActionBadgeClass(stock.action)]">{{ stock.action }}</span>
              </td>
              <td class="text-surface-200/70">{{ stock.rating_from || '—' }}</td>
              <td class="text-surface-200/70">{{ stock.rating_to }}</td>
              <td class="text-surface-200/70">{{ formatPrice(stock.target_from) }}</td>
              <td :class="getTargetChangeClass(stock.target_from, stock.target_to)">
                {{ formatPrice(stock.target_to) }}
              </td>
            </tr>
            <tr v-if="!store.stocks.length">
              <td colspan="8" class="text-center py-12 text-surface-200/40">
                No stocks found. Try a different search or sync data.
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Pagination -->
    <div v-if="store.totalPages > 1" class="flex items-center justify-center gap-2">
      <button
        @click="store.setPage(store.currentPage - 1)"
        :disabled="store.currentPage <= 1"
        class="px-4 py-2 rounded-lg text-sm font-medium transition-all cursor-pointer"
        :class="store.currentPage <= 1
          ? 'bg-surface-800/30 text-surface-200/20 cursor-not-allowed'
          : 'bg-surface-800/50 text-surface-200/70 hover:bg-primary-500/20 hover:text-primary-300'"
      >
        <- Prev
      </button>
      <template v-for="page in store.totalPages" :key="page">
        <button
          v-if="Math.abs(page - store.currentPage) <= 2 || page === 1 || page === store.totalPages"
          @click="store.setPage(page)"
          class="w-10 h-10 rounded-lg text-sm font-medium transition-all cursor-pointer"
          :class="page === store.currentPage
            ? 'bg-primary-500 text-white shadow-lg shadow-primary-500/25'
            : 'bg-surface-800/50 text-surface-200/70 hover:bg-primary-500/20 hover:text-primary-300'"
        >
          {{ page }}
        </button>
        <span
          v-else-if="Math.abs(page - store.currentPage) === 3"
          class="text-surface-200/30"
        >
          ...
        </span>
      </template>
      <button
        @click="store.setPage(store.currentPage + 1)"
        :disabled="store.currentPage >= store.totalPages"
        class="px-4 py-2 rounded-lg text-sm font-medium transition-all cursor-pointer"
        :class="store.currentPage >= store.totalPages
          ? 'bg-surface-800/30 text-surface-200/20 cursor-not-allowed'
          : 'bg-surface-800/50 text-surface-200/70 hover:bg-primary-500/20 hover:text-primary-300'"
      >
        Next ->
      </button>
    </div>

    <!-- Stock Detail Modal -->
    <Teleport to="body">
      <div
        v-if="selectedStock"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        @click.self="selectedStock = null"
      >
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm"></div>
        <div class="relative glass-card p-8 max-w-lg w-full animate-fade-in-up">
          <button
            @click="selectedStock = null"
            class="absolute top-4 right-4 w-8 h-8 rounded-lg bg-white/5 hover:bg-white/10 flex items-center justify-center text-surface-200/50 hover:text-white transition-all cursor-pointer"
          >
            ✕
          </button>
          <div class="mb-6">
            <div class="flex items-center gap-3 mb-2">
              <h2 class="text-3xl font-bold text-white">{{ selectedStock.ticker }}</h2>
              <span :class="['badge', getActionBadgeClass(selectedStock.action)]">{{ selectedStock.action }}</span>
            </div>
            <p class="text-surface-200/50">{{ selectedStock.company }}</p>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div class="bg-surface-800/50 rounded-xl p-4">
              <p class="text-xs text-surface-200/40 mb-1">Brokerage</p>
              <p class="text-white font-medium">{{ selectedStock.brokerage }}</p>
            </div>
            <div class="bg-surface-800/50 rounded-xl p-4">
              <p class="text-xs text-surface-200/40 mb-1">Action</p>
              <p class="text-white font-medium capitalize">{{ selectedStock.action }}</p>
            </div>
            <div class="bg-surface-800/50 rounded-xl p-4">
              <p class="text-xs text-surface-200/40 mb-1">Rating From</p>
              <p class="text-white font-medium">{{ selectedStock.rating_from || '—' }}</p>
            </div>
            <div class="bg-surface-800/50 rounded-xl p-4">
              <p class="text-xs text-surface-200/40 mb-1">Rating To</p>
              <p class="text-white font-medium">{{ selectedStock.rating_to }}</p>
            </div>
            <div class="bg-surface-800/50 rounded-xl p-4">
              <p class="text-xs text-surface-200/40 mb-1">Target From</p>
              <p class="text-white font-medium">{{ formatPrice(selectedStock.target_from) }}</p>
            </div>
            <div class="bg-surface-800/50 rounded-xl p-4">
              <p class="text-xs text-surface-200/40 mb-1">Target To</p>
              <p :class="['font-medium', getTargetChangeClass(selectedStock.target_from, selectedStock.target_to)]">
                {{ formatPrice(selectedStock.target_to) }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>
