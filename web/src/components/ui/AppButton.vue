<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/lib/utils'
import AppSpinner from './AppSpinner.vue'

const props = withDefaults(
  defineProps<{
    variant?: 'primary' | 'secondary' | 'ghost' | 'danger'
    size?: 'sm' | 'md' | 'lg'
    loading?: boolean
    disabled?: boolean
  }>(),
  {
    variant: 'primary',
    size: 'md',
  },
)

const classes = computed(() =>
  cn(
    'inline-flex items-center justify-center gap-2 font-medium rounded-xl transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-accent/50 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer',
    {
      'bg-accent hover:bg-accent-hover text-white shadow-lg shadow-accent-glow hover:shadow-xl hover:shadow-accent/20':
        props.variant === 'primary',
      'bg-white/5 hover:bg-white/10 text-text-primary border border-border':
        props.variant === 'secondary',
      'hover:bg-white/5 text-text-secondary': props.variant === 'ghost',
      'bg-danger/10 hover:bg-danger/20 text-danger': props.variant === 'danger',
    },
    {
      'px-3 py-1.5 text-sm': props.size === 'sm',
      'px-4 py-2.5 text-sm': props.size === 'md',
      'px-6 py-3 text-base': props.size === 'lg',
    },
  ),
)
</script>

<template>
  <button :class="classes" :disabled="disabled || loading">
    <AppSpinner v-if="loading" class="w-4 h-4" />
    <slot />
  </button>
</template>
