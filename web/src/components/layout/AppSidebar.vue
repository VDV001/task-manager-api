<script setup lang="ts">
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { CheckSquare, LayoutDashboard, BarChart3, LogOut } from 'lucide-vue-next'

const route = useRoute()
const auth = useAuthStore()

const navItems = [
  { to: '/', icon: LayoutDashboard, label: 'Tasks' },
  { to: '/stats', icon: BarChart3, label: 'Statistics' },
]
</script>

<template>
  <aside class="w-64 flex-shrink-0 flex flex-col border-r border-border bg-bg-secondary/50 backdrop-blur-sm">
    <!-- Logo -->
    <div class="h-16 flex items-center gap-2.5 px-5 border-b border-border">
      <CheckSquare class="w-6 h-6 text-accent" />
      <span class="text-lg font-bold text-text-primary">TaskFlow</span>
    </div>

    <!-- Nav -->
    <nav class="flex-1 p-3 space-y-1">
      <RouterLink
        v-for="item in navItems"
        :key="item.to"
        :to="item.to"
        :class="[
          'flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all duration-200',
          route.path === item.to
            ? 'bg-accent/10 text-accent'
            : 'text-text-secondary hover:text-text-primary hover:bg-white/5',
        ]"
      >
        <component :is="item.icon" class="w-5 h-5" />
        <span>{{ item.label }}</span>
      </RouterLink>
    </nav>

    <!-- Logout -->
    <div class="p-3 border-t border-border">
      <button
        class="flex items-center gap-3 w-full px-3 py-2.5 rounded-xl text-sm font-medium text-text-secondary hover:text-danger hover:bg-danger/5 transition-all duration-200 cursor-pointer"
        @click="auth.logout(); $router.push('/login')"
      >
        <LogOut class="w-5 h-5" />
        <span>Log out</span>
      </button>
    </div>
  </aside>
</template>
