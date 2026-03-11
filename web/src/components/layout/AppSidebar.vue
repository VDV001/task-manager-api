<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { setLocale, locale } from '@/i18n'
import {
  CheckSquare,
  LayoutDashboard,
  BarChart3,
  LogOut,
  Menu,
  X,
  Languages,
} from 'lucide-vue-next'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const mobileOpen = ref(false)

function handleLogout() {
  auth.logout()
  router.push('/login')
}

function toggleLocale() {
  setLocale(locale.value === 'en' ? 'ru' : 'en')
}

const navItems = [
  { to: '/', icon: LayoutDashboard, labelKey: 'nav.tasks' },
  { to: '/stats', icon: BarChart3, labelKey: 'nav.statistics' },
]
</script>

<template>
  <header
    class="sticky top-0 z-50 w-full bg-bg-primary/80 backdrop-blur-xl border-b border-white/[0.06]"
  >
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex h-16 items-center justify-between">
        <!-- Logo -->
        <div class="flex items-center gap-2.5">
          <CheckSquare class="w-6 h-6 text-accent" />
          <span class="text-lg font-bold text-text-primary tracking-tight">TaskFlow</span>
        </div>

        <!-- Desktop Nav — pill-shaped -->
        <nav class="hidden sm:flex items-center">
          <div
            class="flex items-center gap-1 rounded-full bg-white/[0.04] backdrop-blur-lg border border-white/[0.06] px-1.5 py-1.5"
          >
            <RouterLink v-for="item in navItems" :key="item.to" :to="item.to" class="relative">
              <div
                :class="[
                  'relative flex items-center gap-2 px-4 py-1.5 rounded-full transition-all duration-300 text-sm font-medium',
                  route.path === item.to
                    ? 'text-white'
                    : 'text-text-secondary hover:text-text-primary',
                ]"
              >
                <div
                  v-if="route.path === item.to"
                  class="absolute inset-0 rounded-full bg-gradient-to-r from-accent to-purple-500"
                  style="box-shadow: 0 0 20px rgba(99, 102, 241, 0.4)"
                />
                <div class="relative z-10 flex items-center gap-2">
                  <component :is="item.icon" class="w-4 h-4" />
                  <span>{{ t(item.labelKey) }}</span>
                </div>
              </div>
            </RouterLink>
          </div>
        </nav>

        <!-- Right side -->
        <div class="flex items-center gap-2">
          <!-- Language switcher -->
          <button
            class="flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-medium text-text-secondary hover:text-text-primary hover:bg-white/[0.06] transition-all duration-200 cursor-pointer border border-transparent hover:border-white/[0.08]"
            @click="toggleLocale"
          >
            <Languages class="w-3.5 h-3.5" />
            <span class="uppercase">{{ locale === 'en' ? 'RU' : 'EN' }}</span>
          </button>

          <button
            class="hidden sm:flex items-center gap-2 px-3 py-1.5 rounded-full text-sm font-medium text-text-secondary hover:text-danger hover:bg-danger/5 transition-all duration-200 cursor-pointer"
            @click="handleLogout"
          >
            <LogOut class="w-4 h-4" />
            <span class="hidden md:inline">{{ t('auth.logout') }}</span>
          </button>

          <!-- Mobile menu toggle -->
          <button
            class="sm:hidden p-2 rounded-lg text-text-secondary hover:text-text-primary cursor-pointer"
            @click="mobileOpen = !mobileOpen"
          >
            <X v-if="mobileOpen" class="w-5 h-5" />
            <Menu v-else class="w-5 h-5" />
          </button>
        </div>
      </div>

      <!-- Mobile Nav -->
      <Transition name="slide-down">
        <div v-if="mobileOpen" class="sm:hidden pb-4 space-y-1">
          <RouterLink
            v-for="item in navItems"
            :key="item.to"
            :to="item.to"
            :class="[
              'flex items-center gap-3 px-4 py-2.5 rounded-xl text-sm font-medium transition-all',
              route.path === item.to
                ? 'bg-accent/10 text-accent'
                : 'text-text-secondary hover:text-text-primary hover:bg-white/5',
            ]"
            @click="mobileOpen = false"
          >
            <component :is="item.icon" class="w-5 h-5" />
            {{ t(item.labelKey) }}
          </RouterLink>
          <button
            class="flex items-center gap-3 w-full px-4 py-2.5 rounded-xl text-sm font-medium text-text-secondary hover:text-danger hover:bg-danger/5 transition-all cursor-pointer"
            @click="handleLogout"
          >
            <LogOut class="w-5 h-5" />
            {{ t('auth.logout') }}
          </button>
        </div>
      </Transition>
    </div>
  </header>
</template>

<style scoped>
.slide-down-enter-active,
.slide-down-leave-active {
  transition: all 0.2s ease;
}
.slide-down-enter-from,
.slide-down-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
