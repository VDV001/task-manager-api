<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'

const { t } = useI18n()

defineProps<{
  page: number
  totalPages: number
}>()

defineEmits<{
  'update:page': [page: number]
}>()
</script>

<template>
  <div v-if="totalPages > 1" class="flex items-center justify-center gap-2 mt-6">
    <button
      :disabled="page <= 1"
      class="p-2 rounded-lg hover:bg-white/5 text-text-secondary disabled:opacity-30 transition-colors cursor-pointer disabled:cursor-not-allowed"
      @click="$emit('update:page', page - 1)"
    >
      <ChevronLeft class="w-5 h-5" />
    </button>
    <span class="text-sm text-text-secondary px-3">
      {{ t('pagination.page', { current: page, total: totalPages }) }}
    </span>
    <button
      :disabled="page >= totalPages"
      class="p-2 rounded-lg hover:bg-white/5 text-text-secondary disabled:opacity-30 transition-colors cursor-pointer disabled:cursor-not-allowed"
      @click="$emit('update:page', page + 1)"
    >
      <ChevronRight class="w-5 h-5" />
    </button>
  </div>
</template>
