<script setup lang="ts">
import { cn } from '@/lib/utils'

defineProps<{
  class?: string
  hoverable?: boolean
}>()
</script>

<template>
  <div
    :class="cn(
      'rounded-2xl border border-border bg-bg-card/60 backdrop-blur-sm p-5 transition-all duration-300',
      hoverable && 'hover:border-border-hover hover:bg-bg-card-hover hover:-translate-y-0.5 cursor-pointer',
      $props.class,
    )"
  >
    <slot />
  </div>
</template>
