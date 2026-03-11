<script setup lang="ts">
import { cn } from '@/lib/utils'

defineProps<{
  class?: string
  hoverable?: boolean
}>()
</script>

<template>
  <div
    :class="
      cn(
        'rounded-2xl border border-white/[0.08] bg-white/[0.03] p-5 transition-all duration-300',
        hoverable &&
          'hover:border-white/[0.15] hover:bg-white/[0.05] hover:-translate-y-0.5 hover:shadow-lg hover:shadow-black/20 cursor-pointer',
        $props.class,
      )
    "
  >
    <slot />
  </div>
</template>
