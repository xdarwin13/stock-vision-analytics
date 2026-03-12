<script setup lang="ts">
import { onMounted } from 'vue'
import { useStockStore } from '@/stores/stockStore'

const store = useStockStore()

onMounted(() => {
  store.loadRecommendations(15)
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
  if (score >= 40) return 'from-accent-500 to-accent-400'
  if (score >= 20) return 'from-primary-500 to-primary-400'
  if (score >= 0) return 'from-warning-500 to-warning-400'
  return 'from-danger-500 to-danger-400'
}

function getScoreTextColor(score: number): string {
  if (score >= 40) return 'text-accent-400'
  if (score >= 20) return 'text-primary-400'
  if (score >= 0) return 'text-warning-400'
  return 'text-danger-400'
}

function getScoreLabel(score: number): string {
  if (score >= 40) return 'Strong buy'
  if (score >= 25) return 'Buy'
  if (score >= 15) return 'Moderate'
  if (score >= 0) return 'Hold'
  return 'Avoid'
}

function getTargetChangePercent(from: number, to: number): string {
  if (!from || !to) return '—'
  const change = ((to - from) / from) * 100
  return `${change >= 0 ? '+' : ''}${change.toFixed(1)}%`
}

function getTargetChangeColor(from: number, to: number): string {
  if (!from || !to) return ''
  return to >= from ? 'text-accent-400' : 'text-danger-400'
}

function getRankEmoji(index: number): string {
  if (index === 0) return '🥇'
  if (index === 1) return '🥈'
  if (index === 2) return '🥉'
  return `#${index + 1}`
}
</script>

<template>
  <div>
    <!-- Header -->
    <div class="mb-8 animate-fade-in-up">
      <h1 class="text-3xl font-bold text-white mb-2">🏆 Investment Recommendations</h1>
      <p class="text-surface-200/50">
        AI-powered analysis based on analyst ratings, price targets, and market actions
      </p>
    </div>

    <!-- Algorithm Explanation -->
    <div class="glass-card p-6 mb-8 animate-fade-in-up" style="animation-delay: 0.1s">
      <h3 class="text-lg font-bold text-white mb-3">📐 How We Score</h3>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="flex items-start gap-3">
          <span class="text-xl">📊</span>
          <div>
            <p class="text-sm font-medium text-white">Analyst Action</p>
            <p class="text-xs text-surface-200/40">Upgrades score higher than downgrades (30%)</p>
          </div>
        </div>
        <div class="flex items-start gap-3">
          <span class="text-xl">⬆️</span>
          <div>
            <p class="text-sm font-medium text-white">Rating Change</p>
            <p class="text-xs text-surface-200/40">Positive rating moves boost score (30%)</p>
          </div>
        </div>
        <div class="flex items-start gap-3">
          <span class="text-xl">🎯</span>
          <div>
            <p class="text-sm font-medium text-white">Target Price</p>
            <p class="text-xs text-surface-200/40">Higher price targets are bullish (25%)</p>
          </div>
        </div>
        <div class="flex items-start gap-3">
          <span class="text-xl">⭐</span>
          <div>
            <p class="text-sm font-medium text-white">Rating Strength</p>
            <p class="text-xs text-surface-200/40">"Strong Buy" scores above "Hold" (15%)</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="store.isLoading" class="flex items-center justify-center py-20">
      <div class="flex items-center gap-3 text-surface-200/50">
        <svg class="animate-spin w-5 h-5" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
        </svg>
        <span>Analyzing stocks...</span>
      </div>
    </div>

    <!-- Recommendations List -->
    <div v-else class="space-y-4">
      <div
        v-for="(rec, index) in store.recommendations"
        :key="rec.stock.id"
        class="glass-card p-6 animate-fade-in-up"
        :style="{ animationDelay: `${0.15 + index * 0.05}s` }"
        :class="{ 'pulse-glow border-accent-500/30': index < 3 }"
      >
        <div class="flex flex-col sm:flex-row items-start gap-4">
          <!-- Rank -->
          <div class="flex items-center justify-center w-12 h-12 rounded-xl bg-surface-800/60 text-xl font-bold shrink-0">
            {{ getRankEmoji(index) }}
          </div>

          <!-- Stock Info -->
          <div class="flex-1 min-w-0">
            <div class="flex flex-wrap items-center gap-2 mb-1">
              <span class="text-xl font-bold text-white">{{ rec.stock.ticker }}</span>
              <span :class="['badge', getActionBadgeClass(rec.stock.action)]">{{ rec.stock.action }}</span>
              <span class="badge badge-default">{{ rec.stock.brokerage }}</span>
            </div>
            <p class="text-sm text-surface-200/50 mb-3">{{ rec.stock.company }}</p>

            <!-- Info Grid -->
            <div class="flex flex-wrap gap-x-6 gap-y-2 text-sm">
              <div>
                <span class="text-surface-200/40">Rating: </span>
                <span class="text-white">{{ rec.stock.rating_from || '—' }} → {{ rec.stock.rating_to }}</span>
              </div>
              <div>
                <span class="text-surface-200/40">Target: </span>
                <span class="text-white">{{ formatPrice(rec.stock.target_from) }} → {{ formatPrice(rec.stock.target_to) }}</span>
                <span :class="['ml-1 font-medium', getTargetChangeColor(rec.stock.target_from, rec.stock.target_to)]">
                  ({{ getTargetChangePercent(rec.stock.target_from, rec.stock.target_to) }})
                </span>
              </div>
            </div>

            <!-- Reason -->
            <p class="mt-2 text-xs text-surface-200/40 italic">{{ rec.reason }}</p>
          </div>

          <!-- Score -->
          <div class="flex flex-col items-center sm:items-end shrink-0">
            <p :class="['text-3xl font-bold', getScoreTextColor(rec.score)]">{{ rec.score.toFixed(1) }}</p>
            <p class="text-xs text-surface-200/40">Score</p>
            <div class="mt-2 px-3 py-1 rounded-lg text-xs font-bold bg-gradient-to-r text-white"
              :class="getScoreColor(rec.score)"
            >
              {{ getScoreLabel(rec.score) }}
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-if="!store.recommendations.length" class="glass-card p-12 text-center">
        <p class="text-xl text-surface-200/40 mb-2">No recommendations yet</p>
        <p class="text-sm text-surface-200/30">Add stock data to generate recommendations</p>
      </div>
    </div>
  </div>
</template>
