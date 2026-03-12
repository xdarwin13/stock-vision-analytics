<script setup lang="ts">
import { onMounted } from 'vue'
import { useStockStore } from '@/stores/stockStore'
import { useRouter } from 'vue-router'

const store = useStockStore()
const router = useRouter()

onMounted(async () => {
  await Promise.all([
    store.loadStats(),
    store.loadStocks(),
    store.loadRecommendations(5),
  ])
})

function getActionBadgeClass(action: string): string {
  const normalized = action.toLowerCase()
  if (normalized.includes('upgraded')) return 'badge-upgraded'
  if (normalized.includes('downgraded')) return 'badge-downgraded'
  if (normalized.includes('maintained') || normalized.includes('reiterated')) return 'badge-maintained'
  if (normalized.includes('initiated')) return 'badge-initiated'
  return 'badge-default'
}

function formatPrice(value: number): string {
  if (!value) return '—'
  return `$${value.toFixed(2)}`
}

function getScoreColor(score: number): string {
  if (score >= 40) return 'text-accent-400'
  if (score >= 20) return 'text-primary-400'
  if (score >= 0) return 'text-warning-400'
  return 'text-danger-400'
}
</script>

<template>
  <div>
    <!-- Hero Section -->
    <div class="mb-10 animate-fade-in-up">
      <h1 class="text-4xl sm:text-5xl font-bold mb-3">
        <span class="gradient-text">Stock</span>Vision
      </h1>
      <p class="text-surface-200/60 text-lg max-w-xl">
        Real-time analyst ratings, target prices, and AI-powered investment recommendations.
      </p>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-5 mb-10">
      <div class="glass-card p-6 animate-fade-in-up" style="animation-delay: 0.1s">
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-primary-500/20 to-primary-600/20 flex items-center justify-center text-2xl">
            📈
          </div>
          <div>
            <p class="text-sm text-surface-200/50 font-medium">Total Stocks</p>
            <p class="text-3xl font-bold text-white">{{ store.statsCount }}</p>
          </div>
        </div>
      </div>
      <div class="glass-card p-6 animate-fade-in-up" style="animation-delay: 0.2s">
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-accent-500/20 to-accent-600/20 flex items-center justify-center text-2xl">
            🏆
          </div>
          <div>
            <p class="text-sm text-surface-200/50 font-medium">Top Picks</p>
            <p class="text-3xl font-bold text-white">{{ store.recommendations.length }}</p>
          </div>
        </div>
      </div>
      <div class="glass-card p-6 animate-fade-in-up" style="animation-delay: 0.3s">
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-warning-400/20 to-warning-500/20 flex items-center justify-center text-2xl">
            🔄
          </div>
          <div>
            <p class="text-sm text-surface-200/50 font-medium">Data Source</p>
            <p class="text-lg font-bold text-white">Analyst Ratings</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Top Recommendations -->
    <div class="mb-10 animate-fade-in-up" style="animation-delay: 0.4s">
      <div class="flex items-center justify-between mb-5">
        <h2 class="text-2xl font-bold text-white">🏆 Top Recommendations</h2>
        <button @click="router.push('/recommendations')" class="text-sm text-primary-400 hover:text-primary-300 transition-colors cursor-pointer">
          View all →
        </button>
      </div>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
        <div
          v-for="(rec, index) in store.recommendations.slice(0, 3)"
          :key="rec.stock.id"
          class="glass-card p-6 pulse-glow"
          :style="{ animationDelay: `${index * 0.5}s` }"
        >
          <div class="flex items-start justify-between mb-4">
            <div>
              <div class="flex items-center gap-2 mb-1">
                <span class="text-xl font-bold text-white">{{ rec.stock.ticker }}</span>
                <span :class="['badge', getActionBadgeClass(rec.stock.action)]">
                  {{ rec.stock.action }}
                </span>
              </div>
              <p class="text-sm text-surface-200/50">{{ rec.stock.company }}</p>
            </div>
            <div class="text-right">
              <p :class="['text-2xl font-bold', getScoreColor(rec.score)]">{{ rec.score.toFixed(1) }}</p>
              <p class="text-xs text-surface-200/40">Score</p>
            </div>
          </div>
          <div class="flex items-center gap-4 text-sm">
            <div>
              <p class="text-surface-200/40">Target</p>
              <p class="text-white font-medium">{{ formatPrice(rec.stock.target_from) }} → {{ formatPrice(rec.stock.target_to) }}</p>
            </div>
            <div>
              <p class="text-surface-200/40">Rating</p>
              <p class="text-white font-medium">{{ rec.stock.rating_from || '—' }} → {{ rec.stock.rating_to }}</p>
            </div>
          </div>
          <p class="mt-3 text-xs text-surface-200/40 italic">{{ rec.reason }}</p>
        </div>
      </div>
    </div>

    <!-- Recent Stocks Table -->
    <div class="animate-fade-in-up" style="animation-delay: 0.5s">
      <div class="flex items-center justify-between mb-5">
        <h2 class="text-2xl font-bold text-white">📊 Recent Analyst Actions</h2>
        <button @click="router.push('/stocks')" class="text-sm text-primary-400 hover:text-primary-300 transition-colors cursor-pointer">
          View all →
        </button>
      </div>
      <div class="glass-card overflow-hidden">
        <div class="overflow-x-auto">
          <table class="data-table">
            <thead>
              <tr>
                <th>Ticker</th>
                <th>Company</th>
                <th>Brokerage</th>
                <th>Action</th>
                <th>Rating</th>
                <th>Target</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="stock in store.stocks.slice(0, 8)" :key="stock.id">
                <td class="font-bold text-white">{{ stock.ticker }}</td>
                <td class="text-surface-200/70">{{ stock.company }}</td>
                <td class="text-surface-200/70">{{ stock.brokerage }}</td>
                <td>
                  <span :class="['badge', getActionBadgeClass(stock.action)]">{{ stock.action }}</span>
                </td>
                <td class="text-surface-200/70">{{ stock.rating_from || '—' }} → {{ stock.rating_to }}</td>
                <td class="text-surface-200/70">{{ formatPrice(stock.target_from) }} → {{ formatPrice(stock.target_to) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>
