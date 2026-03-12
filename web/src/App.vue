<script setup lang="ts">
import ErrorBoundary from '@/components/ui/ErrorBoundary.vue'
</script>

<template>
  <ErrorBoundary>
    <RouterView />
  </ErrorBoundary>
</template>
