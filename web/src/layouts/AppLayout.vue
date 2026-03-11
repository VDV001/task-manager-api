<script setup lang="ts">
import { Toaster } from 'vue-sonner'
import AppSidebar from '@/components/layout/AppSidebar.vue'
</script>

<template>
  <div class="min-h-screen bg-bg-primary text-text-primary">
    <AppSidebar />
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <RouterView />
    </main>
    <Toaster
      position="top-right"
      theme="dark"
      :toast-options="{ class: '!bg-bg-card !border-border !text-text-primary' }"
    />
  </div>
</template>
