<script setup lang="ts">
import { Toaster } from 'vue-sonner'
import AppSidebar from '@/components/layout/AppSidebar.vue'
</script>

<template>
  <div class="flex h-screen bg-bg-primary text-text-primary">
    <AppSidebar />
    <main class="flex-1 overflow-y-auto p-8">
      <RouterView />
    </main>
    <Toaster
      position="top-right"
      theme="dark"
      :toast-options="{ class: '!bg-bg-card !border-border !text-text-primary' }"
    />
  </div>
</template>
