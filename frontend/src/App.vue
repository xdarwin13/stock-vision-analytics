<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import { ref } from 'vue'

const navItems = [
  { name: 'Dashboard', path: '/', icon: '📊' },
  { name: 'Stocks', path: '/stocks', icon: '📈' },
  { name: 'Recommendations', path: '/recommendations', icon: '🏆' },
]

const isMobileMenuOpen = ref(false)
</script>

<template>
  <div class="min-h-screen">
    <!-- Navigation -->
    <nav class="fixed top-0 left-0 right-0 z-50 glass-card rounded-none border-x-0 border-t-0">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-16">
          <!-- Logo -->
          <RouterLink to="/" class="flex items-center gap-3 no-underline">
            <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-primary-500 to-primary-700 flex items-center justify-center text-white font-bold text-lg shadow-lg shadow-primary-500/20">
              S
            </div>
            <span class="text-xl font-bold gradient-text">StockVision</span>
          </RouterLink>

          <!-- Desktop Nav -->
          <div class="hidden md:flex items-center gap-1">
            <RouterLink
              v-for="item in navItems"
              :key="item.path"
              :to="item.path"
              class="flex items-center gap-2 px-4 py-2 rounded-xl text-sm font-medium transition-all duration-300 no-underline"
              :class="$route.path === item.path
                ? 'bg-primary-500/20 text-primary-300 shadow-lg shadow-primary-500/10'
                : 'text-surface-200/70 hover:text-white hover:bg-white/5'"
            >
              <span>{{ item.icon }}</span>
              <span>{{ item.name }}</span>
            </RouterLink>
          </div>

          <!-- Mobile menu button -->
          <button
            @click="isMobileMenuOpen = !isMobileMenuOpen"
            class="md:hidden p-2 rounded-lg hover:bg-white/10 text-white"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path v-if="!isMobileMenuOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
              <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>

      <!-- Mobile Nav -->
      <div v-if="isMobileMenuOpen" class="md:hidden border-t border-white/10 pb-4 px-4">
        <RouterLink
          v-for="item in navItems"
          :key="item.path"
          :to="item.path"
          @click="isMobileMenuOpen = false"
          class="flex items-center gap-3 px-4 py-3 mt-1 rounded-xl text-sm font-medium transition-all no-underline"
          :class="$route.path === item.path
            ? 'bg-primary-500/20 text-primary-300'
            : 'text-surface-200/70 hover:text-white hover:bg-white/5'"
        >
          <span>{{ item.icon }}</span>
          <span>{{ item.name }}</span>
        </RouterLink>
      </div>
    </nav>

    <!-- Main Content -->
    <main class="pt-24 pb-12 px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto">
      <RouterView />
    </main>
  </div>
</template>
