<script setup lang="ts">
import { cn } from '@/lib/utils'

defineProps<{
  label?: string
  error?: string
  modelValue?: string
  type?: string
  placeholder?: string
}>()

defineEmits<{
  'update:modelValue': [value: string]
}>()
</script>

<template>
  <div class="flex flex-col gap-1.5">
    <label v-if="label" class="text-sm font-medium text-text-secondary">
      {{ label }}
    </label>
    <input
      :type="type ?? 'text'"
      :value="modelValue"
      :placeholder="placeholder"
      :class="cn(
        'w-full px-4 py-2.5 rounded-xl bg-white/5 border text-text-primary placeholder-text-muted text-sm transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-accent/50',
        error ? 'border-danger focus:ring-danger/50' : 'border-border focus:border-accent',
      )"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    />
    <p v-if="error" class="text-xs text-danger">{{ error }}</p>
  </div>
</template>
