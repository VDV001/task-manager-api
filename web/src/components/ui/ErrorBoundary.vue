<script setup lang="ts">
import { ref, onErrorCaptured } from 'vue'
import { useI18n } from 'vue-i18n'
import { AlertTriangle, RefreshCw } from 'lucide-vue-next'

const { t } = useI18n()
const error = ref<Error | null>(null)
const retryKey = ref(0)

onErrorCaptured((err: Error) => {
  error.value = err
  console.error('[ErrorBoundary]', err)
  return false
})

function retry() {
  error.value = null
  retryKey.value++
}
</script>

<template>
  <slot v-if="!error" :key="retryKey" />
  <div v-else class="flex flex-col items-center justify-center gap-4 py-20 px-6 text-center">
    <div class="w-16 h-16 rounded-2xl bg-danger/10 flex items-center justify-center">
      <AlertTriangle class="w-8 h-8 text-danger" />
    </div>
    <h2 class="text-xl font-semibold text-text-primary">
      {{ t('error.title') }}
    </h2>
    <p class="text-sm text-text-secondary max-w-md">
      {{ t('error.description') }}
    </p>
    <button
      class="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl bg-white/5 hover:bg-white/10 text-text-primary text-sm font-medium border border-border transition-all cursor-pointer"
      @click="retry"
    >
      <RefreshCw class="w-4 h-4" />
      {{ t('error.retry') }}
    </button>
  </div>
</template>
